package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"

	"github.com/lucas-clemente/quic-go"
)

func main() {

	addr := "localhost:8090"

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}
	session, err := quic.DialAddr(addr, tlsConf, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Dialed")

	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println("Stream open")

	fmt.Printf("Client: Sending '%s'\n", "hallo")
	_, err = stream.Write([]byte("hallo"))
	if err != nil {
		panic(err)
	}

	fmt.Println("Waiting for 5 bytes")

	buf := make([]byte, 5)
	_, err = io.ReadFull(stream, buf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Client: Got '%s'\n", buf)

}
