package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/pflag"
)

func main() {
	var timeout time.Duration
	pflag.DurationVar(&timeout, "timeout", time.Second*10, "timeout")
	pflag.Parse()

	args := pflag.Args()
	if len(args) < 2 || len(args[0]) == 0 || len(args[1]) == 0 {
		log.Fatalf("arguments error: %s", args)
	}

	address := net.JoinHostPort(args[0], args[1])

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		log.Fatalf("error while connect: %s", err)
	}

	defer func() {
		_ = client.Close()
	}()

	_, _ = fmt.Fprintf(os.Stderr, "Connected to %s", address)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	go func() {
		defer cancel()
		if err := client.Send(); err != nil {
			log.Fatalf("error while send: %s", err)
		}
	}()

	go func() {
		defer cancel()
		if err := client.Receive(); err != nil {
			log.Fatalf("error while receive: %s", err)
		}
	}()

	<-ctx.Done()
}
