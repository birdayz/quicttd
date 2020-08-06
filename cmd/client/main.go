package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/goiiot/libmqtt"
	"github.com/lucas-clemente/quic-go"
)

func main() {

	addr := "localhost:8090"

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	session, err := quic.DialAddrContext(ctx, addr, tlsConf, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Dialed")

	stream, err := session.OpenStream()
	if err != nil {
		panic(err)
	}

	fmt.Println("Stream open")

	connPacket := libmqtt.ConnPacket{}

	writer := bufio.NewWriter(stream)
	reader := bufio.NewReader(stream)

	fmt.Println("Write")

	err = libmqtt.Encode(&connPacket, writer)
	if err != nil {
		panic(err)
	}

	writer.Flush()

	fmt.Println("Read")

	p, err := libmqtt.Decode(libmqtt.V311, reader)
	if err != nil {
		panic(err)
	}

	stream.Close()

	spew.Dump(p)

}
