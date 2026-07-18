// cmd/setup-s3/main.go — Sets the CORS policy on the S3 bucket so that
// the dashboard (localhost:5173) can do direct browser-to-S3 PUT uploads
// using presigned URLs.
//
// Run once from the backend/ directory:
//   go run ./cmd/setup-s3/main.go
//
// This script reads the same .env as the API server.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	region    := os.Getenv("AWS_REGION")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	bucket    := os.Getenv("AWS_S3_BUCKET")

	if bucket == "" {
		log.Fatal("AWS_S3_BUCKET is not set")
	}

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKey, secretKey, "",
		)),
	)
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	client := s3.NewFromConfig(cfg)

	// Allowed origins: add production domain here when deploying
	allowedOrigins := []string{
		"http://localhost:3000",
		"http://localhost:5173",
	}

	corsRule := types.CORSRule{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{"GET", "PUT", "HEAD"},
		AllowedHeaders: []string{"*"},
		ExposeHeaders:  []string{"ETag", "Content-Length"},
		MaxAgeSeconds:  aws.Int32(3600),
	}

	_, err = client.PutBucketCors(context.Background(), &s3.PutBucketCorsInput{
		Bucket: aws.String(bucket),
		CORSConfiguration: &types.CORSConfiguration{
			CORSRules: []types.CORSRule{corsRule},
		},
	})
	if err != nil {
		log.Fatalf("❌ Failed to set CORS on bucket %q: %v\n\nMake sure your IAM user has s3:PutBucketCors permission.", bucket, err)
	}

	fmt.Printf("✅ CORS policy applied to bucket: %s\n", bucket)
	fmt.Println()
	fmt.Println("Allowed origins:")
	for _, o := range allowedOrigins {
		fmt.Printf("  · %s\n", o)
	}
	fmt.Println()
	fmt.Println("Allowed methods : GET, PUT, HEAD")
	fmt.Println("Allowed headers : * (all)")
	fmt.Println("Exposed headers : ETag, Content-Length")
	fmt.Println("MaxAgeSeconds   : 3600")
	fmt.Println()
	fmt.Println("You can verify in AWS Console → S3 → your bucket → Permissions → CORS")
}
