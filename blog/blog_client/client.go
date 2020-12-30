package main

import (
	"context"
	"fmt"
	"io"
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

	///*
	//.########..########....###....########.....########..##........#######...######..
	//.##.....##.##.........##.##...##.....##....##.....##.##.......##.....##.##....##.
	//.##.....##.##........##...##..##.....##....##.....##.##.......##.....##.##.......
	//.########..######...##.....##.##.....##....########..##.......##.....##.##...####
	//.##...##...##.......#########.##.....##....##.....##.##.......##.....##.##....##.
	//.##....##..##.......##.....##.##.....##....##.....##.##.......##.....##.##....##.
	//.##.....##.########.##.....##.########.....########..########..#######...######..
	//*/
	//Will throw an error
	_, err2 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "1jh3hh38887"})

	if err2 != nil {
		fmt.Printf("Error Happened while reading %v \n", err2)
	}

	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogID}
	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), readBlogReq)
	if readBlogErr != nil {
		fmt.Println("This error is supposed to happen")
		fmt.Printf("Error Happened while reading %v \n", readBlogErr)
	}
	fmt.Printf("Blog was read: %v \n", readBlogRes)
	// doUnary(c)

	///*
	//.##.....##.########..########.....###....########.########.########..##........#######...######..
	//.##.....##.##.....##.##.....##...##.##......##....##.......##.....##.##.......##.....##.##....##.
	//.##.....##.##.....##.##.....##..##...##.....##....##.......##.....##.##.......##.....##.##.......
	//.##.....##.########..##.....##.##.....##....##....######...########..##.......##.....##.##...####
	//.##.....##.##........##.....##.#########....##....##.......##.....##.##.......##.....##.##....##.
	//.##.....##.##........##.....##.##.....##....##....##.......##.....##.##.......##.....##.##....##.
	//..#######..##........########..##.....##....##....########.########..########..#######...######..
	//*/
	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Changed Author",
		Title:    "My first Blog (edited)",
		Content:  "Content of the first blog, with some awesome additions",
	}
	updateRes, updateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})

	if updateErr != nil {
		fmt.Println("error happened while updating ", updateErr)
	}
	fmt.Println("Blog was updated: ", updateRes)

	/*
		.########..########.##.......########.########.########....########..##........#######...######..
		.##.....##.##.......##.......##..........##....##..........##.....##.##.......##.....##.##....##.
		.##.....##.##.......##.......##..........##....##..........##.....##.##.......##.....##.##.......
		.##.....##.######...##.......######......##....######......########..##.......##.....##.##...####
		.##.....##.##.......##.......##..........##....##..........##.....##.##.......##.....##.##....##.
		.##.....##.##.......##.......##..........##....##..........##.....##.##.......##.....##.##....##.
		.########..########.########.########....##....########....########..########..#######...######..
	*/
	deleteRes, deleteErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogID})
	if deleteErr != nil {
		fmt.Println("error happened while deleting ", deleteErr)
	}
	fmt.Println("Blog was deleted: ", deleteRes)

	/*
		.##.......####..######..########....########..##........#######...######....######.
		.##........##..##....##....##.......##.....##.##.......##.....##.##....##..##....##
		.##........##..##..........##.......##.....##.##.......##.....##.##........##......
		.##........##...######.....##.......########..##.......##.....##.##...####..######.
		.##........##........##....##.......##.....##.##.......##.....##.##....##........##
		.##........##..##....##....##.......##.....##.##.......##.....##.##....##..##....##
		.########.####..######.....##.......########..########..#######...######....######.
	*/
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})

	if err != nil {
		log.Fatalf("error while calling ListBlog RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something Happened %v", err)
		}
		fmt.Println(res.GetBlog())

	}
}
