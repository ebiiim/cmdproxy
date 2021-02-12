package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ebiiim/cmdproxy"

	"golang.org/x/crypto/acme/autocert"
)

func main() {
	secret := flag.String("secret", "", "Secret")
	isTLS := flag.Bool("tls", false, "Use TLS")
	host := flag.String("host", "127.0.0.1:12345", "TLS: example.com, non-TLS: 127.0.0.1:12345")
	flag.Parse()

	if *secret == "" || *host == "" {
		flag.Usage()
		os.Exit(1)
	}

	s := cmdproxy.NewServer(*secret)

	mux := http.NewServeMux()
	mux.HandleFunc(s.Path(), s.Run)

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if *isTLS {
			log.Println(http.Serve(autocert.NewListener(*host), mux))
			sigCh <- syscall.SIGTERM
		} else {
			log.Println(http.ListenAndServe(*host, mux))
			sigCh <- syscall.SIGTERM
		}
	}()
	sig := <-sigCh
	fmt.Printf("Signal <%s> received. Shutting down...\n", sig)
}
