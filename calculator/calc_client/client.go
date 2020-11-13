package main

import (
	"context"
	"log"

	"github.com/grpcproj/calculator/calcpb"
	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect", err)
	}
	defer cc.Close()

	c := calcpb.NewSumServiceClient(cc)

	request := &calcpb.SumRequest{
		Sum: &calcpb.Sum{
			FirstNumber: 1,
			LastNumber:  1,
		},
	}

	response, err := c.Add(context.Background(), request)
	if err != nil {
		log.Fatalf("error while calling greet RPC: %v", err)
	}

	log.Printf("the respose was: %v", response.Result)
}
