package main

import (
	"context"
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
		log.Fatalf("Could not connect %v", err)
	}

	defer cc.Close()
	//create the client
	c := blogpb.NewBlogServiceClient(cc)

	blog := &blogpb.Blog{
		AuthorId: "Stephane",
		Title:    "My first Blog",
		Content:  "Content of the first blog",
	}
	fmt.Println("Creating blog")
	createdBlogResponse, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: %v \n", err)
	}
	fmt.Printf("Blog has been created: %v \n", createdBlogResponse)
	blogID := createdBlogResponse.GetBlog().GetId()

	//Read Blog
	//Will throw an error
	_, err2 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "1jh3hh38887"})

	if err2 != nil {
		fmt.Printf("Error Happened while reading %v \n", err2)
	}

	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogID}
	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), readBlogReq)
	if readBlogErr != nil {
		fmt.Printf("Error Happened while reading %v \n", readBlogErr)
	}
	fmt.Printf("Blog was read: %v \n", readBlogRes)
	// doUnary(c)
}
