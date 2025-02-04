package main

import (
	"context"
	"log"
	"time"

	pb "ldtgrpc01/proto" // replace with your module name

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewCalculatorClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Make the gRPC call
	r, err := c.Add(ctx, &pb.AddRequest{Num1: 5, Num2: 3})
	if err != nil {
		log.Fatalf("could not calculate: %v", err)
	}
	log.Printf("Result: %d", r.GetResult())
}
