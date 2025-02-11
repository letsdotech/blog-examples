package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	pb "ldtgrpc01/proto" // replace with your module name
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct {
	pb.UnimplementedCalculatorServer
}

func main() {
	// Load server certificate and private key
	cert, err := tls.LoadX509KeyPair("../certs/server-cert.pem", "../certs/server-key.pem")
	if err != nil {
		log.Fatalf("failed to load server certificates: %v", err)
	}

	// Create a certificate pool and add the client's CA certificate
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("../certs/ca-cert.pem")
	if err != nil {
		log.Fatalf("failed to read ca certificate: %v", err)
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatal("failed to append client certs")
	}

	// Create the TLS credentials
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
		MinVersion:   tls.VersionTLS12,
	})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterCalculatorServer(s, &server{})
	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Add method implementation
func (s *server) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	result := req.Num1 + req.Num2
	return &pb.AddResponse{Result: result}, nil
}
