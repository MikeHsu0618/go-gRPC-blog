package repository

import (
	"context"
	"io"
	"log"

	pb "go-grpc-blog/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Repository interface {
	CreateBlog(ctx context.Context, data *pb.Blog) string
	GetBlogById(ctx context.Context, id string) *pb.Blog
	UpdateBlog(ctx context.Context, data *pb.Blog)
	ListBlog(ctx context.Context)
	DeleteBlog(ctx context.Context, id string)
}

type repository struct {
	conn pb.BlogServiceClient
}

func (repo *repository) DeleteBlog(ctx context.Context, id string) {
	log.Println("DeleteBlog was invoked")
	_, err := repo.conn.DeleteBlog(ctx, &pb.BlogId{Id: id})
	if err != nil {
		log.Fatalf("Unexpected error: %v \n", err)
	}
	log.Println("Success DeleteBlog")
}

func (repo *repository) ListBlog(ctx context.Context) {
	log.Println("ListBlog was invoked")
	stream, err := repo.conn.ListBlog(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("Error happened while ListBlog : %v \n", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println("Unexpected Error", err)
		}
		log.Println(res)
	}
}

func NewRepository(conn pb.BlogServiceClient) Repository {
	return &repository{conn: conn}
}

func (repo *repository) CreateBlog(ctx context.Context, data *pb.Blog) string {
	log.Println("CreateBlog was invoked")
	res, err := repo.conn.CreateBlog(ctx, data)
	if err != nil {
		log.Fatalf("Unexpected error: %v \n", err)
	}

	return res.Id
}

func (repo *repository) GetBlogById(ctx context.Context, id string) *pb.Blog {
	log.Println("ReadBlog was invoked")
	req := &pb.BlogId{Id: id}
	res, err := repo.conn.ReadBlog(ctx, req)

	if err != nil {
		log.Fatalf("Error happened while reading: %v", err)
	}

	log.Printf("Blog was read: %v", res)
	return res
}

func (repo *repository) UpdateBlog(ctx context.Context, data *pb.Blog) {
	log.Println("UpdateBlog was invoked")

	_, err := repo.conn.UpdateBlog(ctx, data)

	if err != nil {
		log.Fatalf("Error happened while updating: %v", err)
	}

	log.Println("Success Update", data)
}
