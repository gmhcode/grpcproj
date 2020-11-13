package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/grpcproj/greet/greetpb"

	"google.golang.org/grpc"
)

type server struct{}

//When the server receives the greeting, it will respond with this VV
func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v \n", req)
	// prints
	// Greet function was invoked with greeting:{first_name:"Greg" last_name:"Hughes"}
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil

}

func main() {
	fmt.Println("Hello world")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {

		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	//Gets the "Greet" functions ready
	greetpb.RegisterGreetServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)

	}

}
