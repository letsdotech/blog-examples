package main

import (
	"context"
	"io"
	pb "ldtgrpc03/proto"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	// Load client certificates
	creds, err := credentials.NewClientTLSFromFile(
		"certs/server-cert.pem", // Server certificate
		"localhost",             // Server name (must match the certificate's CN)
	)
	if err != nil {
		log.Fatalf("Failed to load credentials: %v", err)
	}

	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewCalculatorClient(conn)

	// Unary RPC
	// unaryExample(client)

	// // Server Streaming RPC
	// serverStreamingExample(client)

	// // Client Streaming RPC
	// clientStreamingExample(client)

	// // Bidirectional Streaming RPC
	bidirectionalStreamingExample(client)
}

func unaryExample(client pb.CalculatorClient) {
	ctx := context.Background()
	resp, err := client.Add(ctx, &pb.AddRequest{Num1: 10, Num2: 20})
	if err != nil {
		log.Fatalf("could not add: %v", err)
	}
	log.Printf("Sum: %d", resp.Result)
}

func serverStreamingExample(client pb.CalculatorClient) {
	ctx := context.Background()
	stream, err := client.GenerateNumbers(ctx, &pb.GenerateRequest{Limit: 5})
	if err != nil {
		log.Fatalf("error calling GenerateNumbers: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error receiving: %v", err)
		}
		log.Printf("Received number: %d", resp.Number)
	}
}

func clientStreamingExample(client pb.CalculatorClient) {
	ctx := context.Background()
	stream, err := client.ComputeAverage(ctx)
	if err != nil {
		log.Fatalf("error calling ComputeAverage: %v", err)
	}

	numbers := []int64{1, 2, 3, 4, 5}
	for _, num := range numbers {
		if err := stream.Send(&pb.NumberRequest{Number: num}); err != nil {
			log.Fatalf("error sending: %v", err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error receiving response: %v", err)
	}
	log.Printf("Average: %.2f", resp.Result)
}

func bidirectionalStreamingExample(client pb.CalculatorClient) {
	ctx := context.Background()
	stream, err := client.ProcessNumbers(ctx)
	if err != nil {
		log.Fatalf("error calling ProcessNumbers: %v", err)
	}

	waitc := make(chan struct{})

	// Send numbers
	go func() {
		numbers := []int64{1, 2, 3, 4, 5}
		for _, num := range numbers {
			if err := stream.Send(&pb.NumberRequest{Number: num}); err != nil {
				log.Fatalf("error sending: %v", err)
			}
			time.Sleep(500 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// Receive processed numbers
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("error receiving: %v", err)
			}
			log.Printf("Received processed number: %d", resp.Number)
		}
	}()

	<-waitc
}
