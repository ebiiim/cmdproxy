package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ebiiim/cmdproxy"
)

func main() {
	secret := flag.String("secret", "", "Secret")
	url := flag.String("url", "http://localhost:12345/", "URL")
	cmd := flag.String("cmd", "echo helloworld", "Command")
	flag.Parse()
	if *secret == "" || *url == "" || *cmd == "" {
		flag.Usage()
		os.Exit(1)
	}
	client := cmdproxy.NewClient(*url, *secret)

	res, err := client.Run(strings.Split((*cmd), " "), 5*time.Second)
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}
	fmt.Printf("Error: %+v\n", res.Error)
	fmt.Printf("ExitCode: %d\n", res.ExitCode)
	fmt.Printf("Stdout: %s\n", res.Stdout)
	fmt.Printf("Stderr: %s\n", res.Stderr)
}
