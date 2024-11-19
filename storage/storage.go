package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

type StorageClient struct {
	client     *storage.Client
	projectId  string
	bucketName string
	path       string
}

func Init(path string) *StorageClient {
	if os.Getenv("CLOUD_SERVICE_ACCOUNT_KEY") != "" {
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "keys/"+os.Getenv("CLOUD_SERVICE_ACCOUNT_KEY"))
	}

	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return &StorageClient{
		client:     client,
		projectId:  os.Getenv("CLOUD_PROJECT_ID"),
		bucketName: os.Getenv("CLOUD_BUCKET_NAME"),
		path:       path,
	}
}

func (c *StorageClient) UploadFile(file multipart.File, object string) (string, error) {
	ctx := context.Background()

	ctx, timeout := context.WithTimeout(ctx, time.Second*50)
	defer timeout()

	wc := c.client.Bucket(c.bucketName).Object(c.path + object).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}

	return "https://storage.googleapis.com/" + c.bucketName + "/" + c.path + object, nil
}
