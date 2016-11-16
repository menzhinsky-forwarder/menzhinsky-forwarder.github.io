package main

import (
	"flag"
	"io"
	"log"
	"net"
)

var (
	flagNet         string
	flagListenAddr  string
	flagForwardAddr string
)

func main() {
	flag.StringVar(&flagNet, "net", "tcp", "network type")
	flag.StringVar(&flagListenAddr, "listen", ":8080", "address to listen to")
	flag.StringVar(&flagForwardAddr, "forward", ":8080", "addres to forward traffic to")
	flag.Parse()

	listener, err := net.Listen(flagNet, flagListenAddr)
	if err != nil {
		log.Fatalf("listen err: %v", err)
	}

	for {
		var conn net.Conn
		conn, err = listener.Accept()
		if err != nil {
			log.Fatalf("accept err: %v", err)
		}

		go forward(conn)
	}
}

func forward(src net.Conn) {
	dst, err := net.Dial(flagNet, flagForwardAddr)
	if err != nil {
		src.Close()
		log.Printf("dial err=%v", err)
		return
	}

	go func() {
		_, err := io.Copy(src, dst)
		if err != nil {
			src.Close()
			log.Printf("copy dst->src err: %v", err)
		}
	}()

	go func() {
		_, err := io.Copy(dst, src)
		if err != nil {
			src.Close()
			log.Printf("copy src->dst err: %v", err)
		}
	}()
}
