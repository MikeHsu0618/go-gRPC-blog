package main

import (
	"fmt"
	"log"

	"go-grpc-blog/client/repository"
	"go-grpc-blog/client/service"
	pb "go-grpc-blog/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr string = "0.0.0.0:50051"

func main() {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Couldn't connect to client: %v\n", err)
	}

	defer conn.Close()
	c := pb.NewBlogServiceClient(conn)
	repo := repository.NewRepository(c)
	svc := service.NewService(repo)

	// createBlog
	id := svc.CreateBlog("testId", "title", "content")

	// readBlog
	res := svc.ReadBlog(id)
	fmt.Println("Success ReadBlog", res)

	// updateBlog
	svc.UpdateBlog(&pb.Blog{
		Id:       id,
		AuthorId: "updateAAA",
		Title:    "updateTitle",
		Content:  "updateContent",
	})

	// listBlog
	svc.ListBlog()

	// deleteBlog
	svc.DeleteBlog(id)
}
