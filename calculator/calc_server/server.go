package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"github.com/grpcproj/calculator/calcpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct{}

func (*server) Add(ctx context.Context, req *calcpb.SumRequest) (*calcpb.SumResponse, error) {
	fmt.Printf("Add function was invoked with %v \n", req)
	// prints
	// Greet function was invoked with greeting:{first_name:"Greg" last_name:"Hughes"}
	firstNumber := req.GetSum().GetFirstNumber()
	lastNumber := req.GetSum().GetLastNumber()
	result := firstNumber + lastNumber
	res := &calcpb.SumResponse{
		Result: result,
	}
	return res, nil

}

func (*server) PrimeNumberDecomposition(req *calcpb.PrimeNumberDecompositionRequest, stream calcpb.SumService_PrimeNumberDecompositionServer) error {
	fmt.Printf("Received PrimeNumberDecomposition RPC: %v\n", req)
	number := req.GetNumber()
	divisor := int64(2)
	for number > 1 {
		//we send the divisor back to the client
		if number%divisor == 0 {
			stream.Send(&calcpb.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			})
			number = number / divisor
		} else {
			divisor++
			fmt.Printf("Divisor has increased to  %v\n", divisor)
		}
	}
	return nil
}
func (*server) ComputeAverage(stream calcpb.SumService_ComputeAverageServer) error {
	fmt.Printf("Received ComputeAverage RPC: \n")
	sum := int32(0)
	count := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			average := float64(sum) / float64(count)
			return stream.SendAndClose(&calcpb.ComputeAverageResponse{
				Average: average,
			})
		}

		if err != nil {
			log.Fatalf("error while reading client stream: %v", err)
		}
		sum += req.GetNumber()
		count++
	}
}

func (*server) FindMaximum(stream calcpb.SumService_FindMaximumServer) error {
	fmt.Printf("Received FindMaximum rpc")
	maximum := int32(0)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error while reading client stream: %v", err)
			return err
		}
		number := req.GetNumber()
		///anytime the number is higher than the maximum, we update the maximum and send it back.
		if number > maximum {
			maximum = number

			sendErr := stream.Send(&calcpb.FindMaximumResponse{
				Maximum: maximum,
			})
			if sendErr != nil {
				log.Fatalf("error while sending data to client: %v", err)
				return sendErr
			}
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *calcpb.SquareRootRequest) (*calcpb.SquareRootResponse, error) {
	fmt.Println("Received SquareRoot RPC")
	number := req.GetNumber()

	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received negative number: %v", number),
		)
	}
	return &calcpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func main() {
	fmt.Println("Running Server")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {

		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	//Gets the "Greet" functions ready
	calcpb.RegisterSumServiceServer(s, &server{})
	//Register reflection service on the GRPC Server
	//use "evans" commant line tool to find the parameters for this.. which is basically postman for GRPC
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)

	}

}
