package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
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

// deleteS3Object deletes the specified object from the S3 bucket
func deleteS3Object(ctx context.Context, client *s3.Client, bucketName, objectKey string) error {
	input := &s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	}

	_, err := client.DeleteObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete object: %v", err)
	}

	// Notification of successful deletion
	fmt.Printf("Successfully deleted object: %s\n", objectKey)

	return nil
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

	if len(result.Contents) == 0 {
		fmt.Println("No objects found in the bucket.")
		os.Exit(0) // Exit if no objects are found
	}

	// Show file size in MB
	for _, obj := range result.Contents {
		sizeInMB := float64(*obj.Size) / (1024 * 1024)
		fmt.Printf("Key: %s, Size: %.2f MB\n", *obj.Key, sizeInMB)
	}

	return nil
}

// downloadS3Object downloads the specified object from the S3 bucket
func downloadS3Object(ctx context.Context, client *s3.Client, bucketName, objectKey string) error {
	output, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	})
	if err != nil {
		return fmt.Errorf("failed to download object: %v", err)
	}
	defer output.Body.Close()

	// Create a file to save the downloaded object
	file, err := os.Create(objectKey)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Copy the object data to the file
	_, err = io.Copy(file, output.Body)
	if err != nil {
		return fmt.Errorf("failed to save object to file: %v", err)
	}

	// Notification of successful download
	fmt.Printf("Successfully downloaded object: %s\n", objectKey)

	return nil
}

// uploadS3Object uploads the specified file to the S3 bucket
func uploadS3Object(ctx context.Context, client *s3.Client, bucketName, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	input := &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    aws.String(filepath.Base(filePath)), // Use only the file name as the object key
		Body:   file,
	}

	_, err = client.PutObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to upload object: %v", err)
	}

	// Notification of successful upload
	fmt.Printf("Successfully uploaded object: %s\n", filepath.Base(filePath))

	return nil
}

// menu handles user input for deleting, downloading, or uploading objects or exiting
func menu(ctx context.Context, client *s3.Client, bucketName string) {
	for {
		err := listS3BucketContents(ctx, client, bucketName)
		if err != nil {
			log.Fatalf("Failed to list bucket contents: %v", err)
		}

		var action string
		fmt.Println("Enter 'delete' to delete an object, 'download' to download an object, 'upload' to upload a file, or 'exit' to exit:")
		fmt.Scanln(&action)

		switch action {
		case "exit":
			return
		case "delete":
			var objectKey string
			fmt.Println("Enter the object key to delete:")
			fmt.Scanln(&objectKey)
			err := deleteS3Object(ctx, client, bucketName, objectKey)
			if err != nil {
				log.Printf("Error deleting object: %v", err)
			}
		case "download":
			var objectKey string
			fmt.Println("Enter the object key to download:")
			fmt.Scanln(&objectKey)
			err := downloadS3Object(ctx, client, bucketName, objectKey)
			if err != nil {
				log.Printf("Error downloading object: %v", err)
			}
		case "upload":
			var filePath string
			fmt.Println("Enter the file path to upload:")
			fmt.Scanln(&filePath)
			err := uploadS3Object(ctx, client, bucketName, filePath)
			if err != nil {
				log.Printf("Error uploading object: %v", err)
			}
		default:
			fmt.Println("Invalid command. Please try again.")
		}
	}
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

	menu(ctx, client, bucketName)
}
