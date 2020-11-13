package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/grpcproj/calculator/calcpb"
	"google.golang.org/grpc"
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
func main() {
	fmt.Println("Listening for Sums")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {

		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	//Gets the "Greet" functions ready
	calcpb.RegisterSumServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)

	}

}
