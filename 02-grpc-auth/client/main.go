package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"time"

	pb "ldtgrpc01/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

func verifyConnection(state tls.ConnectionState) {
	log.Printf("=== TLS Connection Details ===")
	log.Printf("Version: %x", state.Version)
	log.Printf("CipherSuite: %s", tls.CipherSuiteName(state.CipherSuite))
	log.Printf("HandshakeComplete: %t", state.HandshakeComplete)
	log.Printf("Server Name: %s", state.ServerName)

	for i, cert := range state.PeerCertificates {
		log.Printf("Peer Certificate [%d]:", i)
		log.Printf("  Subject: %s", cert.Subject)
		log.Printf("  Issuer: %s", cert.Issuer)
		log.Printf("  Valid from: %s", cert.NotBefore)
		log.Printf("  Valid until: %s", cert.NotAfter)
	}
}

func main() {
	// Load client certificate and private key
	cert, err := tls.LoadX509KeyPair("../certs/client-cert.pem", "../certs/client-key.pem")
	if err != nil {
		log.Fatalf("failed to load client certificates: %v", err)
	}

	// Create a certificate pool and add the CA certificate
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("../certs/ca-cert.pem")
	if err != nil {
		log.Fatalf("failed to read ca certificate: %v", err)
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatal("failed to append ca certs")
	}

	// Create TLS config with verification callbacks
	tlsConfig := &tls.Config{
		ServerName:   "localhost",
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
		MinVersion:   tls.VersionTLS12,
		VerifyConnection: func(cs tls.ConnectionState) error {
			verifyConnection(cs)
			return nil
		},
	}

	// Create the TLS credentials
	creds := credentials.NewTLS(tlsConfig)

	// Create connection with verification
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewCalculatorClient(conn)

	// Make multiple requests to test connection
	for i := 0; i < 3; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// Get peer info before making the call
		if p, ok := peer.FromContext(ctx); ok {
			if tlsInfo, ok := p.AuthInfo.(credentials.TLSInfo); ok {
				log.Printf("Connection using TLS version: %x", tlsInfo.State.Version)
			}
		}

		r, err := c.Add(ctx, &pb.AddRequest{Num1: int32(i), Num2: 20})
		if err != nil {
			log.Fatalf("could not calculate: %v", err)
		}
		log.Printf("Result: %d", r.GetResult())
		time.Sleep(1 * time.Second)
	}
}
