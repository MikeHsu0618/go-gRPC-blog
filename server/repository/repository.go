package repository

import (
	"context"
	"errors"

	pb "go-grpc-blog/proto"
	"go-grpc-blog/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	InsertData(ctx context.Context, data *model.BlogItem) (*mongo.InsertOneResult, error)
	GetInsertedID(res *mongo.InsertOneResult) primitive.ObjectID
	GetOid(id string) (primitive.ObjectID, error)
	GetBlogByOid(ctx context.Context, oid primitive.ObjectID) (*pb.Blog, error)
	UpdateBlog(ctx context.Context, data *model.BlogItem, oid primitive.ObjectID) (*mongo.UpdateResult, error)
	GetListBlog(ctx context.Context) (*mongo.Cursor, error)
	DeleteBlog(ctx context.Context, oid primitive.ObjectID) error
	DocumentToBlog(data *model.BlogItem) *pb.Blog
}

type repository struct {
	collection *mongo.Collection
}

func (repo *repository) DeleteBlog(ctx context.Context, oid primitive.ObjectID) error {
	res, err := repo.collection.DeleteOne(ctx, bson.M{"_id": oid})

	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("blog was not found")
	}

	return nil
}

func NewRepository(collection *mongo.Collection) Repository {
	return &repository{collection: collection}
}

func (repo *repository) InsertData(ctx context.Context, data *model.BlogItem) (*mongo.InsertOneResult, error) {
	res, err := repo.collection.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (repo *repository) GetInsertedID(res *mongo.InsertOneResult) primitive.ObjectID {
	return res.InsertedID.(primitive.ObjectID)
}

func (repo *repository) GetOid(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}

func (repo *repository) GetBlogByOid(ctx context.Context, oid primitive.ObjectID) (*pb.Blog, error) {
	data := &model.BlogItem{}
	filter := bson.M{"_id": oid}
	res := repo.collection.FindOne(ctx, filter)
	if err := res.Decode(data); err != nil {
		return nil, err
	}

	return repo.DocumentToBlog(data), nil
}

func (repo *repository) UpdateBlog(ctx context.Context, data *model.BlogItem, oid primitive.ObjectID) (*mongo.UpdateResult, error) {
	var res, err = repo.collection.UpdateOne(
		ctx,
		bson.M{"_id": oid},
		bson.M{"$set": data},
	)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (repo *repository) GetListBlog(ctx context.Context) (*mongo.Cursor, error) {
	cur, err := repo.collection.Find(ctx, primitive.D{{}})

	if err != nil {
		return nil, err
	}
	return cur, nil
}

func (repo *repository) DocumentToBlog(data *model.BlogItem) *pb.Blog {
	return &pb.Blog{
		Id:       data.ID.Hex(),
		AuthorId: data.AuthorId,
		Title:    data.Title,
		Content:  data.Content,
	}
}
