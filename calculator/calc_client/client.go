package main

import (
	"context"
	"fmt"
	"io"
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
	// doUnary(c)
	doServerStreaming(c)

}

func doUnary(c calcpb.SumServiceClient) {
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

//doStream
func doServerStreaming(c calcpb.SumServiceClient) {
	fmt.Println("Starting to do a PrimeDecomposition Server Streaming RPC: ")
	req := &calcpb.PrimeNumberDecompositionRequest{

		Number: 123145,
	}

	stream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PrimeDecomposition RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something Happened %v", err)
		}
		fmt.Println(res.GetPrimeFactor())
	}

	// log.Printf("the respose was: %v", response.Result)
}
