package service

import (
	"context"
	"fmt"
	"log"

	pb "go-grpc-blog/proto"
	"go-grpc-blog/server/model"
	"go-grpc-blog/server/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type service struct {
	pb.UnimplementedBlogServiceServer
	repo repository.Repository
}

func NewService(repo repository.Repository) pb.BlogServiceServer {
	return &service{pb.UnimplementedBlogServiceServer{}, repo}
}

func (s *service) CreateBlog(ctx context.Context, in *pb.Blog) (*pb.BlogId, error) {
	log.Printf("CreateBlog was invoked with %v", in)
	res, err := s.repo.InsertData(ctx, &model.BlogItem{
		AuthorId: in.AuthorId,
		Title:    in.Title,
		Content:  in.Content,
	})

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal Error : %v \n", err),
		)
	}

	oid := s.repo.GetInsertedID(res)

	return &pb.BlogId{Id: oid.Hex()}, nil
}

func (s *service) ReadBlog(ctx context.Context, in *pb.BlogId) (*pb.Blog, error) {
	log.Printf("ReadBlog was invoked with : %v \n\n", in)
	oid, err := s.repo.GetOid(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Cant Parse ID : %v \n", err))
	}
	res, err := s.repo.GetBlogByOid(ctx, oid)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Cant find blog with Id : %v \n", in.Id))
	}

	return res, nil
}

func (s *service) UpdateBlog(ctx context.Context, in *pb.Blog) (*emptypb.Empty, error) {
	log.Printf("UpdateBlog waws invoked with %v \n", in)

	oid, err := primitive.ObjectIDFromHex(in.Id)

	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cant Parse ID"),
		)
	}
	data := &model.BlogItem{
		AuthorId: in.AuthorId,
		Title:    in.Title,
		Content:  in.Content,
	}
	res, err := s.repo.UpdateBlog(ctx, data, oid)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"Cant Update",
		)
	}

	if res.MatchedCount == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			"Cant find blog with Id",
		)
	}

	return &emptypb.Empty{}, nil
}

func (s *service) DeleteBlog(ctx context.Context, in *pb.BlogId) (*emptypb.Empty, error) {
	log.Printf("DeleteBlog was invoked")

	oid, err := s.repo.GetOid(in.Id)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			err.Error(),
		)
	}

	err = s.repo.DeleteBlog(ctx, oid)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			err.Error(),
		)
	}
	return &emptypb.Empty{}, nil
}

func (s *service) ListBlog(empty *emptypb.Empty, stream pb.BlogService_ListBlogServer) error {
	log.Printf("ListBlogs was invoked")
	cur, err := s.repo.GetListBlog(context.Background())
	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unexpect Internal Error: %v", err),
		)
	}

	for cur.Next(context.Background()) {
		data := &model.BlogItem{}
		err := cur.Decode(data)
		if err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Error while decoding data from MongoDB : %v", err),
			)
		}

		stream.Send(s.repo.DocumentToBlog(data))
	}

	if err = cur.Err(); err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unexpect Internal Error: %v", err),
		)
	}

	return nil
}
