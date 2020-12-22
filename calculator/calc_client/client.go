package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/grpcproj/calculator/calcpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	// doBiDiStreaming(c)
	doErrorUnary(c)
}

func doUnary(c calcpb.SumServiceClient) {
	fmt.Println("Starting to do a Sum Unary RPC...")
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

func doErrorUnary(c calcpb.SumServiceClient) {
	fmt.Println("Starting to do a QuareRoot Unary RPC...")

	// correct call
	fmt.Println("Doing Correct call")
	doErrorCall(c, 10)
	// error call
	fmt.Println("Doing error call")
	doErrorCall(c, -2)

	//we want it to print
	// Result of SquareRoot of 10: 3.1622776601683795Received negative number: -2
	// InvalidArgument
	// We probably sent a negative number
	// Result of SquareRoot of -2: 0
}

func doErrorCall(c calcpb.SumServiceClient, num int32) {
	res, err := c.SquareRoot(context.Background(), &calcpb.SquareRootRequest{Number: num})

	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			//actual error from grpc (good..we created it)
			fmt.Printf("Error message from server: %v \n", respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number")
				return
			}
		} else {
			log.Fatalf("Big Error we didnt create calling SquareRoot: %v", err)
			return
		}
	}

	fmt.Printf("Result of SquareRoot of %v: %v \n", num, res.GetNumberRoot())
}
