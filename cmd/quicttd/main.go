package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"

	"bufio"

	"github.com/davecgh/go-spew/spew"
	"github.com/goiiot/libmqtt"
	quic "github.com/lucas-clemente/quic-go"
)

func handleSession(session quic.Session) error {
	for {
		stream, err := session.AcceptStream(context.Background())
		if err != nil {
			return err
		}

		go handleStream(stream)
	}

}

func handleStream(stream quic.Stream) error {
	packet, err := libmqtt.Decode(libmqtt.V311, bufio.NewReader(stream))
	if err != nil {
		return err
	}

	connAck := libmqtt.ConnAckPacket{}

	w := bufio.NewWriter(stream)

	err = libmqtt.Encode(&connAck, w)
	if err != nil {
		return err
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	spew.Dump(packet)

	stream.Close()

	return nil

}

func main() {
	addr := "localhost:8090"
	lnr, err := quic.ListenAddr(addr, generateTLSConfig(), nil)
	if err != nil {
		panic(err)
	}

	for {
		fmt.Println("w8 accept")
		session, err := lnr.Accept(context.Background())
		if err != nil {
			panic(err)
		}

		go func() {
			fmt.Println("Waiting for sess close")
			<-session.Context().Done()
			fmt.Println("SESS DONE!")
		}()

		go handleSession(session)
	}
}

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}

	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo-example"},
	}
}
