package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/aws"
	"mime/multipart"
	"log"
   )

var uploader *manager.Uploader

// SetUpS3Uploader initializes the S3 uploader
func SetUpS3Uploader() {
	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("error loading AWS config: %v", err)
	}

	// Create an S3 client
	client := s3.NewFromConfig(cfg)

	// Create an S3 uploader
	uploader = manager.NewUploader(client)
}

// UploadFile uploads a file to the specified S3 bucket
func UploadFile(file *multipart.FileHeader) (*manager.UploadOutput, error) {
	// Ensure the uploader is set up
	if uploader == nil {
		SetUpS3Uploader()
	}

	//Open file
	f, err := file.Open()
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	// Upload the file
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("social-app-declan"),
		Key:    aws.String(file.Filename),
		Body:   f,
		ACL: "public-read",

	})
	if err != nil {
		return nil, err
	}

	return result, nil
}