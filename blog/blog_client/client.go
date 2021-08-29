package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/carloscontrerasruiz/grpc-course/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("Hello Im the client")

	opts := grpc.WithInsecure()

	//opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:50051", opts)

	if err != nil {
		log.Fatalf("Could not conect: %v", err)
	}

	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	//doUnary(c)
	//readBlog(c, "61292c45db83f05710e389d4")
	//updateBlog(c, "61292c45db83f05710e389d4")
	//deleteBlog(c, "61292c45db83f05710e389d4")
	listBlog(c)
}

func doUnary(c blogpb.BlogServiceClient) {

	blog := &blogpb.Blog{
		AuthorId: "Carlos",
		Title:    "Mi primer blog",
		Content:  "Este es mi primer blog grpv en go",
	}

	fmt.Println("Creating a blog entry")

	response, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
		Blog: blog,
	})
	if err != nil {
		log.Fatalf("Could not create blog: %v", err)
	}

	fmt.Printf("Blog has been created %v", response)

}

func readBlog(c blogpb.BlogServiceClient, blogId string) {
	_, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: "sdsdsdsdsdsdsd",
	})
	if err != nil {
		log.Printf("Could not read blog: %v", err)
	}

	res, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: blogId,
	})
	if err != nil {
		log.Fatalf("Could not read blog: %v", err)
	}

	fmt.Printf("Blog has been read %v", res)

}

func updateBlog(c blogpb.BlogServiceClient, id string) {
	newBlog := &blogpb.Blog{
		Id:       id,
		AuthorId: "NEw author",
		Title:    "new title",
		Content:  "new content",
	}

	updateRes, err := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{
		Blog: newBlog,
	})
	if err != nil {
		log.Printf("Could not update blog: %v", err)
	}

	fmt.Printf("Blog has been updated %v", updateRes)
}

func deleteBlog(c blogpb.BlogServiceClient, id string) {

	res, err := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{
		BlogId: id,
	})
	if err != nil {
		log.Printf("Could not update blog: %v", err)
	}

	fmt.Printf("Blog has deleted updated %v", res)
}

func listBlog(c blogpb.BlogServiceClient) {

	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Printf("Could not list blog1: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not list blog2: %v", err)
		}

		fmt.Println(res.GetBlog())
	}

	fmt.Println("Blogs listed")
}
