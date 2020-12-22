package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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
	// doServerStreaming(c)
	// doClientStreaming(c)
	doBiDiStreaming(c)
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

func doClientStreaming(c calcpb.SumServiceClient) {
	fmt.Println("Starting to do a ComputeAverage Client Streaming RPC...")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while opening stream: %v", err)
	}

	numbers := []int32{3, 5, 9, 54, 23}

	for _, number := range numbers {
		fmt.Printf("Sending number: %v\n", number)
		stream.Send(&calcpb.ComputeAverageRequest{
			Number: number,
		})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response: %v", err)
	}

	fmt.Printf("The Average is: %v\n", res.GetAverage())

}

func doBiDiStreaming(c calcpb.SumServiceClient) {
	fmt.Println("Starting to do a FindMaximum BiDi Streaming RPC...")

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error While opening stream and calling FindMaximum: %v", err)
	}
	waitc := make(chan struct{})

	//send go routine
	go func() {
		numbers := []int32{4, 7, 2, 19, 4, 6, 32}
		for _, number := range numbers {
			fmt.Printf("Sending number %v\n", number)
			stream.Send(&calcpb.FindMaximumRequest{
				Number: number,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		//indicated client is done streaming numbers
		stream.CloseSend()
	}()
	//receive go routine
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break

			}
			if err != nil {
				log.Fatalf("Problem while reading server stream: %v", err)
				break
			}
			maximum := res.GetMaximum()
			fmt.Printf("Received a new Maximum of:... %v \n", maximum)
		}
		//closes channel
		close(waitc)
	}()
	<-waitc
}
