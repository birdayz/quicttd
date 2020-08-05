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

	quic "github.com/lucas-clemente/quic-go"
)

func main() {
	addr := "localhost:8090"
	lnr, err := quic.ListenAddr(":8080", generateTLSConfig(), nil)
	if err != nil {
		panic(err)
	}

	session, err := lnr.Accept(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println("Accepted Session")

	stream, err := session.AcceptStream(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println("Accepted Stream")

	n, err := stream.Write([]byte("hallo"))
	if err != nil {
		fmt.Println("Wrote bytes", n)
	}

	stream.Close()

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
