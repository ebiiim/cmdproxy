package cmdproxy

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

var (
	ServerLogging = true
)

var (
	ErrClientEncode = errors.New("ErrClientEncode")
	ErrClientPOST   = errors.New("ErrClientPOST")
	ErrClientDecode = errors.New("ErrClientDecode")
)

type req struct {
	Secret  []byte        `json:"secret"`
	Cmd     []string      `json:"cmd"`
	Timeout time.Duration `json:"timeout"`
}

type Result struct {
	Error    string `json:"error"`
	ExitCode int    `json:"exitcode"`
	Stdout   []byte `json:"stdout"`
	Stderr   []byte `json:"stderr"`
}

type Client struct {
	secret []byte
	url    string
	client *http.Client
}

func NewClient(server, secret string) *Client {
	client := &http.Client{}
	c := &Client{
		secret: getSecretHash(secret),
		url:    strings.TrimRight(server, "/") + "/" + getPathHash(secret),
		client: client,
	}
	return c
}

func (c *Client) Run(cmd []string, timeout time.Duration) (*Result, error) {
	r := &req{
		Secret:  c.secret,
		Cmd:     cmd,
		Timeout: timeout,
	}
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(r); err != nil {
		return nil, fmt.Errorf("%w (%v)", ErrClientEncode, err)
	}
	resp, err := c.client.Post(c.url, "application/json", &b)
	if err != nil {
		return nil, fmt.Errorf("%w (%v)", ErrClientPOST, err)
	}
	defer resp.Body.Close()
	var result Result
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("%w (%v)", ErrClientDecode, err)
	}
	return &result, nil
}

type Server struct {
	secret []byte
	path   string
}

func NewServer(secret string) *Server {
	s := &Server{
		secret: getSecretHash(secret),
		path:   "/" + getPathHash(secret),
	}
	return s
}

func (s *Server) Path() string {
	return s.path
}

func (s *Server) Run(w http.ResponseWriter, r *http.Request) {
	var req req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		if ServerLogging {
			log.Printf("BadRequest Method=%v RemoteAddr=%v UA=%v", r.Method, r.RemoteAddr, r.UserAgent())
		}
		writeResult(w, http.StatusBadRequest, err, -1, nil, nil)
		return
	}
	if bytes.Compare(req.Secret, s.secret) != 0 {
		if ServerLogging {
			log.Printf("Unauthorized Method=%v RemoteAddr=%v UA=%v", r.Method, r.RemoteAddr, r.UserAgent())
		}
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ctx, cFn := context.WithTimeout(context.Background(), req.Timeout)
	defer cFn()
	if ServerLogging {
		log.Printf("OK Method=%v RemoteAddr=%v UA=%v", r.Method, r.RemoteAddr, r.UserAgent())
		log.Printf("%s\n", strings.Join(req.Cmd, " "))
	}
	cmd := exec.CommandContext(ctx, req.Cmd[0], req.Cmd[1:]...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	writeResult(w, 200, err, cmd.ProcessState.ExitCode(), stdout.Bytes(), stderr.Bytes())
}

func writeResult(w http.ResponseWriter, code int, resErr error, resCode int, resStdout, resStderr []byte) {
	sErr := ""
	if resErr != nil {
		sErr = resErr.Error()
	}
	res := &Result{
		Error:    sErr,
		ExitCode: resCode,
		Stdout:   resStdout,
		Stderr:   resStderr,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(res)
}

func getSecretHash(s string) []byte {
	return getHash(s, "cmdproxyCMDPROXYcmdproxy")
}

func getPathHash(s string) string {
	return hex.EncodeToString(getHash(s, "CMDPROXYcmdproxyCMDPROXY"))
}

func getHash(s, salt string) []byte {
	r := sha256.Sum256([]byte(s + salt))
	return r[:]
}
