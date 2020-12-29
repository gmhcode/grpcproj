package main

import (
	"fmt"
	"log"

	"github.com/grpcproj/blog/blogpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog Client")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	// cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect", err)
	}

	defer cc.Close()
	//create the client
	c := blogpb.NewBlogServiceClient(cc)
	// doUnary(c)
}
