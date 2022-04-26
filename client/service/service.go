package service

import (
	"context"
	"log"

	"go-grpc-blog/client/repository"
	pb "go-grpc-blog/proto"
)

var ctx = context.Background()

type Service interface {
	CreateBlog(authorId string, title string, content string) string
	ReadBlog(id string) *pb.Blog
	UpdateBlog(blog *pb.Blog)
	ListBlog()
	DeleteBlog(id string)
}

type service struct {
	repo repository.Repository
}

func (s *service) DeleteBlog(id string) {
	s.repo.DeleteBlog(ctx, id)
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) UpdateBlog(blog *pb.Blog) {
	s.repo.UpdateBlog(ctx, blog)
}

func (s *service) CreateBlog(authorId string, title string, content string) string {
	blogId := s.repo.CreateBlog(ctx, &pb.Blog{
		AuthorId: authorId,
		Title:    title,
		Content:  content,
	})
	log.Printf("Blog has been created: %s \n", blogId)
	return blogId
}

func (s *service) ReadBlog(id string) *pb.Blog {
	res := s.repo.GetBlogById(ctx, id)
	return res
}

func (s *service) ListBlog() {
	s.repo.ListBlog(ctx)
}
