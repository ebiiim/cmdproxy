# cmdproxy

Run any commands on remote hosts via HTTP(S).

## Usage

1. Install and run `cmdprx` on the remote host.

```sh
go get "github.com/ebiiim/cmdproxy/..."
go install "github.com/ebiiim/cmdproxy/cmd/cmdprx"

cmdprx -host="10.0.0.100:12345" -secret="abc"
```

2. Install and run `cmdprx-cli` locally, then you can get results from the remote host.

```sh
go get "github.com/ebiiim/cmdproxy/..."
go install "github.com/ebiiim/cmdproxy/cmd/cmdprx-cli"

cmdprx-cli -url="http://10.0.0.100:12345" -secret="abc" -cmd="echo hello world"
```

```
Error: 
ExitCode: 0
Stdout: hello world

Stderr: 

```

## Planned Features

### Security

- [ ] Out-of-the-box TLS support (Let's Encrypt)
- [ ] Client certificate authentication
- [ ] Source IP filter
- [ ] Request header filter
- [ ] Fixed source port

### Feature

- [ ] Stdin support (but no way to support streams, use gRPC?)
- [ ] Encode errors (Base64 encoded gob?)

### Usability

- [ ] Log enhancement
- [ ] Documentation
