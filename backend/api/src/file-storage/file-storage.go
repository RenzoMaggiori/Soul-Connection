package filestorage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Open(connectionString string) (*mongo.Client, context.Context, error) {
	ctx := context.TODO()
	filestorage, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))

	if err != nil {
		return nil, nil, err
	}
	if err := filestorage.Ping(ctx, nil); err != nil {
		return nil, nil, err
	}
	return filestorage, ctx, nil
}

func ConnectionString() string {
	mongo_user := os.Getenv("MONGO_USER")
	mongo_password := os.Getenv("MONGO_PASSWORD")
	mongo_host := os.Getenv("MONGO_HOST")
	mongo_port := os.Getenv("MONGO_PORT")

	return fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin", mongo_user, mongo_password, mongo_host, mongo_port)
}

func Upload(bucket *gridfs.Bucket, r io.Reader, filename string) (*primitive.ObjectID, error) {
	fileBytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	uploadStream, err := bucket.OpenUploadStream(filename)
	if err != nil {
		return nil, err
	}
	defer uploadStream.Close()

	_, err = uploadStream.Write(fileBytes)
	if err != nil {
		return nil, err
	}

	fileObjectId, ok := uploadStream.FileID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("Could not convert upload stream file id to object id")
	}
	return &fileObjectId, nil
}

func DownloadByName(bucket *gridfs.Bucket, filename string) ([]byte, error) {
	downloadStream, err := bucket.OpenDownloadStreamByName(filename)
	if err != nil {
		return nil, err
	}
	defer downloadStream.Close()

	fileContent, err := io.ReadAll(downloadStream)
	if err != nil {
		return nil, err
	}
	return fileContent, nil
}

func DownloadById(bucket *gridfs.Bucket, fileId primitive.ObjectID) ([]byte, error) {
	downloadStream, err := bucket.OpenDownloadStream(fileId)
	if err != nil {
		return nil, err
	}
	defer downloadStream.Close()

	fileContent, err := io.ReadAll(downloadStream)
	if err != nil {
		return nil, err
	}
	return fileContent, nil
}

func Delete(bucket *gridfs.Bucket, fileId primitive.ObjectID) error {
	err := bucket.Delete(fileId)
	if err != nil {
		return err
	}
	return nil
}
