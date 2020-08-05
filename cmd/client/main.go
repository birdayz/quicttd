package main

import (
	"bufio"
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
	session, err := quic.DialAddr(addr, tlsConf, nil)
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

	spew.Dump(p)

}
