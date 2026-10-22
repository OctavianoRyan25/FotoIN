package infra

import (
	"context"
	"fmt"
	"io"
	"time"

	"fotoin/backend/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	client    *s3.Client
	bucket    string
	region    string
	endpoint  string
	presigner *s3.PresignClient
}

func NewS3Client(ctx context.Context, cfg config.S3Config) (*S3Client, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKey, cfg.SecretKey, "",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("load AWS config: %w", err)
	}

	opts := []func(*s3.Options){}

	// Support for MinIO / S3-compatible storage
	if cfg.Endpoint != "" {
		opts = append(opts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
			o.UsePathStyle = true
		})
	}

	client := s3.NewFromConfig(awsCfg, opts...)

	// Automatically ensure bucket exists when running with custom endpoint (MinIO / Localstack)
	if cfg.Endpoint != "" && cfg.Bucket != "" {
		_, err := client.HeadBucket(ctx, &s3.HeadBucketInput{
			Bucket: aws.String(cfg.Bucket),
		})
		if err != nil {
			_, _ = client.CreateBucket(ctx, &s3.CreateBucketInput{
				Bucket: aws.String(cfg.Bucket),
			})
		}
	}

	// PresignClient: we do NOT set presign options globally here.
	// Per-request options are passed in each method below.
	presigner := s3.NewPresignClient(client)

	return &S3Client{
		client:    client,
		bucket:    cfg.Bucket,
		region:    cfg.Region,
		endpoint:  cfg.Endpoint,
		presigner: presigner,
	}, nil
}

// GetPresignedUploadURL generates a presigned URL for direct browser-to-S3 PUT upload.
//
// IMPORTANT — S3 CORS setup required (run once in AWS console or with aws-cli):
//
//	aws s3api put-bucket-cors --bucket YOUR_BUCKET --cors-configuration '{
//	  "CORSRules": [{
//	    "AllowedHeaders": ["*"],
//	    "AllowedMethods": ["PUT", "GET", "HEAD"],
//	    "AllowedOrigins": ["http://localhost:5173", "https://yourdomain.com"],
//	    "ExposeHeaders":  ["ETag"],
//	    "MaxAgeSeconds":  3600
//	  }]
//	}'
//
// The contentType parameter is included in the presigned signature.
// The browser PUT request MUST send the exact same Content-Type header,
// otherwise AWS returns 403 SignatureDoesNotMatch.
func (s *S3Client) GetPresignedUploadURL(key, contentType string, ttl time.Duration) (string, error) {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}

	req, err := s.presigner.PresignPutObject(context.Background(), input,
		s3.WithPresignExpires(ttl),
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned upload URL: %w", err)
	}
	return req.URL, nil
}

// GetPresignedDownloadURL generates a presigned URL for private object download
func (s *S3Client) GetPresignedDownloadURL(key string, ttl time.Duration) (string, error) {
	req, err := s.presigner.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(ttl))
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned download URL: %w", err)
	}
	return req.URL, nil
}

// UploadObject uploads a file directly from the backend (server-side) to S3.
// Use this as a fallback proxy when the browser cannot PUT directly to S3
// due to missing CORS policy on the bucket.
func (s *S3Client) UploadObject(key, contentType string, body io.Reader, size int64) error {
	_, err := s.client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		ContentType:   aws.String(contentType),
		Body:          body,
		ContentLength: aws.Int64(size),
	})
	if err != nil {
		return fmt.Errorf("s3 upload failed: %w", err)
	}
	return nil
}

// DeleteObject removes an object from S3
func (s *S3Client) DeleteObject(key string) error {
	_, err := s.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}

// GetPublicURL returns the public URL for watermarked/preview images
func (s *S3Client) GetPublicURL(key string) string {
	if s.endpoint != "" {
		return fmt.Sprintf("%s/%s/%s", s.endpoint, s.bucket, key)
	}
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, key)
}
