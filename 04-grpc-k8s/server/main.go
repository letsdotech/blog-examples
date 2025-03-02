package main

import (
	"context"
	pb "ldtgrpc04/proto"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct {
	pb.UnimplementedCalculatorServer
}

// Unary RPC - Server
func (s *server) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	println("Add method called. ")
	result := req.Num1 + req.Num2
	return &pb.AddResponse{Result: result}, nil
}

// Server Streaming RPC
func (s *server) GenerateNumbers(req *pb.GenerateRequest, stream pb.Calculator_GenerateNumbersServer) error {
	println("GenerateNumbers method called. ")
	for i := int64(0); i < req.Limit; i++ {
		if err := stream.Send(&pb.NumberResponse{Number: i}); err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

// Client Streaming RPC
func (s *server) ComputeAverage(stream pb.Calculator_ComputeAverageServer) error {
	var sum int64
	var count int64
	println("ComputeAverage method called. ")

	for {
		req, err := stream.Recv()
		if err != nil {
			return stream.SendAndClose(&pb.AverageResponse{
				Result: float64(sum) / float64(count),
			})
		}
		println("Received number: ", req.Number)
		sum += req.Number
		count++
	}
}

// Bidirectional Streaming RPC
func (s *server) ProcessNumbers(stream pb.Calculator_ProcessNumbersServer) error {
	println("ProcessNumbers method called. ")
	for {
		req, err := stream.Recv()
		if err != nil {
			return nil
		}

		// Process the number (multiply by 2) and send it back
		result := req.Number * 2
		if err := stream.Send(&pb.NumberResponse{Number: result}); err != nil {
			return err
		}
		println("Received number: ", req.Number)
	}
}

func main() {
	// Load server certificate and private key
	creds, err := credentials.NewServerTLSFromFile(
		"certs/server-cert.pem", // Server certificate
		"certs/server-key.pem",  // Server private key
	)
	if err != nil {
		log.Fatalf("Failed to load certificates: %v", err)
	}

	// Create gRPC server with TLS credentials
	s := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterCalculatorServer(s, &server{})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
