package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	if flag.NArg() != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s [--timeout=10s] host port\n", os.Args[0])
		os.Exit(1)
	}

	host := flag.Arg(0)
	port := flag.Arg(1)

	address := net.JoinHostPort(host, port)

	conn, err := net.DialTimeout("tcp", address, *timeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connection error: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	done := make(chan struct{}, 2)

	// STDIN -> socket
	go func() {
		_, _ = io.Copy(conn, os.Stdin)
		_ = conn.Close()

		done <- struct{}{}
	}()

	// socket -> STDOUT
	go func() {
		_, _ = io.Copy(os.Stdout, conn)

		done <- struct{}{}
	}()

	<-done
}
