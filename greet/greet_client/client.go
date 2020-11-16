package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/grpcproj/stream/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("hello i am a client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect", err)
	}

	defer cc.Close()
	//create the client
	c := greetpb.NewGreetServiceClient(cc)
	// doUnary(c)
	doServerStreaming(c)

}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Streaming Client RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Greg",
			LastName:  "Hughes",
		},
	}
	//contex is things like wait timers and stuff, background i guess is just like a nil context
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling greet RPC: %v", err)
	}
	log.Printf("Respnse from Greet: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Greg",
			LastName:  "Hughes",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream", err)
		}
		log.Printf("response from GreetManyTimes: %v", msg.GetResult())
	}

}
