package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/grpcproj/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func generateSSLOptions() grpc.DialOption {
	certFile := "ssl/ca.crt" //Certificate Authority Trust Certificate
	creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
	if sslErr != nil {
		log.Fatalf("Error While loading ca trust certificate: %v", sslErr)
		return nil
	}
	//SSL credentials
	opts := grpc.WithTransportCredentials(creds)
	return opts
}

func main() {
	fmt.Println("hello i am a client")

	tls := false
	opts := grpc.WithInsecure()
	//when tls is true, it will use SSL stuff
	if tls {
		opts = generateSSLOptions()
	}

	cc, err := grpc.Dial("localhost:50051", opts)
	// cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect", err)
	}

	defer cc.Close()
	//create the client
	c := greetpb.NewGreetServiceClient(cc)
	doUnary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	// doBiDiStreaming(c)

	// fmt.Println("giving Server 5 sec to respond.. should respond in 3 sec")
	// doUnaryWithDealine(c, 5*time.Second) //should complete
	// fmt.Println("giving Server 1 sec to respond.. should respond in 3 sec")
	// doUnaryWithDealine(c, 1*time.Second) //should timeout

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

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")
	//Dont neet to pass request because it is a stream
	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Greg",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Greg 2",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Greg 3",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Greg 4",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Greg 5",
			},
		},
	}
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LONGGREET RPC: %v", err)
	}
	// we interate over our slice and send each message individually
	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from LongGReet RPC: %v", err)
	}
	fmt.Printf("LongGreet Response: %v\n", res)
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
func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	// we create a stream by invoking the client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while creating stream", err)
		return
	}

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Greg",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Greg 2",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Greg 3",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Greg 4",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Greg 5",
			},
		},
	}
	waitc := make(chan struct{})
	go func() {
		// function to send a bunch of messages
		for _, req := range requests {
			fmt.Printf("Sending Message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	go func() {
		// function to receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while Receiving %v", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	<-waitc
	//we send a bunch of messages to the client (go routine)

	//we receive a bunch of messages from the client (go routine)

	// block until everything is done.

}

func doUnaryWithDealine(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Starting to do a Streaming Client RPC...")
	req := &greetpb.GreetWithDealineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Greg",
			LastName:  "Hughes",
		},
	}
	//contex is things like wait timers and stuff, background i guess is just like a nil context
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetWithDealine(ctx, req)

	if err != nil {

		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline was exceeded")
			} else {
				fmt.Printf("unexpected Error: %v \n", statusErr)
			}
		} else {
			log.Fatalf("error whil calline GreetWithDealine RPC: %v \n", err)
		}
		return
	}
	log.Printf("Respnse from Greet: %v", res.Result)
}
