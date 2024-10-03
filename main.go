package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

// initS3Client initializes and returns an S3 client
func initS3Client(ctx context.Context) (*s3.Client, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	// Create a custom AWS config loader
	customProvider := config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
		os.Getenv("AWS_ACCESS_KEY_ID"),
		os.Getenv("AWS_SECRET_ACCESS_KEY"),
		"",
	))

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("eu-west-1"),
		customProvider,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS econfiguration: %v", err)
	}
	return s3.NewFromConfig(cfg), nil
}

// listS3BucketContents lists the contents of the specified S3 bucket
func listS3BucketContents(ctx context.Context, client *s3.Client, bucketName string) error {
	input := &s3.ListObjectsV2Input{
		Bucket: &bucketName,
	}

	result, err := client.ListObjectsV2(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to list objects: %v", err)
	}

	for _, obj := range result.Contents {
		fmt.Printf("Key: %s, Size: %d bytes\n", *obj.Key, obj.Size)
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <bucket-name>")
		os.Exit(1)
	}

	bucketName := os.Args[1]
	ctx := context.Background()

	client, err := initS3Client(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize S3 client: %v", err)
	}

	err = listS3BucketContents(ctx, client, bucketName)
	if err != nil {
		log.Fatalf("Failed to list bucket contents: %v", err)
	}
}
