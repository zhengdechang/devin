package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"path"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioS3Service provides an implementation of StorageService using the MinIO Go SDK.
// It can be used to connect to MinIO servers as well as other S3-compatible services
// like AWS S3 and Aliyun OSS (when S3 compatibility is enabled).
type MinioS3Service struct {
	client   *minio.Client
	region   string // Optional, but often needed for AWS S3
	// BaseURL for constructing preview URLs. Can be CDN or direct S3 endpoint.
	// e.g., "https://my-cdn.com" or "https://s3.us-east-1.amazonaws.com"
	// If using pre-signed URLs, this might not be directly used in GetPreviewURL
	// but is good for general configuration.
	previewBaseURL string
}

// NewMinioS3Service creates a new MinioS3Service instance.
// endpoint: The S3-compatible endpoint (e.g., "s3.amazonaws.com", "oss-cn-hangzhou.aliyuncs.com", "localhost:9000").
// accessKeyID: The access key.
// secretAccessKey: The secret key.
// useSSL: Boolean indicating whether to use HTTPS.
// region: The AWS region (e.g., "us-east-1"). Can be an empty string if not applicable (e.g., for local MinIO).
// previewBaseURL: The base URL for constructing public preview links. Example: "https://your-bucket.s3.region.amazonaws.com" or "https://cdn.yourdomain.com"
func NewMinioS3Service(endpoint, accessKeyID, secretAccessKey, region, previewBaseURL string, useSSL bool) (*MinioS3Service, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
		Region: region, // Set region for S3 services
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	// Optional: Ping the server to ensure connectivity, though MinIO client doesn't have a direct ping.
	// A lightweight operation like ListBuckets (even if it fails due to permissions) can verify endpoint and credentials.
	// For now, we assume the client initializes correctly if no error is returned by minio.New.

	return &MinioS3Service{
		client:         minioClient,
		region:         region,
		previewBaseURL: previewBaseURL,
	}, nil
}

// Upload uploads a file to the specified bucket with the given object key.
func (s *MinioS3Service) Upload(bucketName string, objectKey string, reader io.Reader, objectSize int64, contentType string) (string, error) {
	ctx := context.Background()

	// Ensure bucket exists, create if not.
	// For production, you might want more sophisticated bucket management or assume buckets are pre-created.
	exists, err := s.client.BucketExists(ctx, bucketName)
	if err != nil {
		return "", fmt.Errorf("failed to check if bucket '%s' exists: %w", bucketName, err)
	}
	if !exists {
		// MakeBucketOptions can take a region, defaults to client's region if empty.
		err = s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: s.region})
		if err != nil {
			return "", fmt.Errorf("failed to create bucket '%s': %w", bucketName, err)
		}
		// fmt.Printf("Successfully created bucket %s\n", bucketName) // For debugging
	}

	uploadInfo, err := s.client.PutObject(ctx, bucketName, objectKey, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
		// Other options like UserMetadata, SSE, etc., can be added here.
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload object '%s' to bucket '%s': %w", objectKey, bucketName, err)
	}

	// uploadInfo.Key is the object key, which we already have.
	// uploadInfo.Location can sometimes be useful but might not be the final public URL.
	return uploadInfo.Key, nil
}

// Download retrieves a file from the specified bucket with the given object key.
func (s *MinioS3Service) Download(bucketName string, objectKey string) (io.ReadCloser, error) {
	ctx := context.Background()
	object, err := s.client.GetObject(ctx, bucketName, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to download object '%s' from bucket '%s': %w", objectKey, bucketName, err)
	}
	return object, nil
}

// GetPreviewURL generates a URL suitable for previewing an object.
// This implementation constructs a URL using the configured previewBaseURL, bucketName, and objectKey.
// It assumes that the objects are publicly readable or that the previewBaseURL points to a CDN/proxy
// that handles access. For private objects, pre-signed URLs would be needed.
func (s *MinioS3Service) GetPreviewURL(bucketName string, objectKey string, configuredBaseURL string) (string, error) {
	// Use the configuredBaseURL passed during factory creation for consistency.
	// If a more dynamic baseURL per call is needed, the interface/method signature might need adjustment.
	if configuredBaseURL == "" {
		// Fallback if no specific baseURL was provided via factory config, try to construct a sensible default
		// This is a very basic S3-like path.
		// For MinIO, if it's self-hosted, the endpoint in minioClient might be usable.
		// For AWS S3, it's typically https://<bucket-name>.s3.<region>.amazonaws.com/<object-key>
		// For Aliyun OSS, it's typically https://<bucket-name>.<endpoint>/<object-key>
		// This part can be complex due to different URL structures.
		// A robust solution might require the user to always provide a correct previewBaseURL.
		return "", fmt.Errorf("previewBaseURL is not configured for the storage service; cannot generate preview URL without it or pre-signing")
	}

	// Ensure objectKey does not have leading slashes if baseURL already ends with one
	// or if path.Join is used, which handles this.

	base, err := url.Parse(configuredBaseURL)
	if err != nil {
		return "", fmt.Errorf("invalid previewBaseURL '%s': %w", configuredBaseURL, err)
	}

	// If the configuredBaseURL is for the S3 service itself (e.g. https://s3.region.amazonaws.com)
	// then the bucket name should be part of the path.
	// If configuredBaseURL is a bucket-specific URL (e.g. https://mybucket.s3.region.amazonaws.com)
	// then bucketName should not be added to the path again.
	// This logic depends heavily on how `configuredBaseURL` is intended to be used.
	// For this implementation, we assume `configuredBaseURL` is a general prefix,
	// and we append `bucketName` and `objectKey`.
	// Example: configuredBaseURL = "https://my-cdn.com", results in "https://my-cdn.com/bucketName/objectKey"
	// Example: configuredBaseURL = "https://s3.region.amazonaws.com", results in "https://s3.region.amazonaws.com/bucketName/objectKey"

	// A common pattern for S3 direct URLs (non-CDN) when `configuredBaseURL` is the S3 endpoint (e.g. `https://s3.region.amazonaws.com`)
	// is `<configuredBaseURL>/<bucketName>/<objectKey>`.
	// If `configuredBaseURL` is a custom domain or CDN that maps to a bucket (e.g. `https://my-cdn.com` -> `my-bucket`),
	// then the URL might be `<configuredBaseURL>/<objectKey>`.
	// The current `StorageService` interface takes `bucketName` and `objectKey` separately for `GetPreviewURL`,
	// and the factory config also has a `BaseURL`. We will use the factory's `BaseURL`.

	// Let's assume the `configuredBaseURL` from the factory is the true base, and path.Join will construct the rest.
	// This means if `configuredBaseURL` is "https://my-cdn.com/my-bucket-path", then the result will be
	// "https://my-cdn.com/my-bucket-path/objectKey" if `bucketName` is not part of the path segments here.
	// Or, if the `configuredBaseURL` is "https://s3.region.amazonaws.com", then we'd want
	// "https://s3.region.amazonaws.com/bucketName/objectKey".

	// The plan stated: "The `GetPreviewURL` will construct a URL by combining the provided `baseURL`, `bucketName`, and `objectKey`."
	// The `baseURL` here refers to the one from the factory config.

	// Path components to join with the base URL.
	// We will join bucketName and objectKey as path segments.
	// This might not be correct for all S3 URL styles (e.g., virtual-hosted style vs. path-style).
	// For simplicity, we use path-style.
	// Virtual-hosted: https://bucket-name.s3.region-code.amazonaws.com/key-name
	// Path-style:    https://s3.region-code.amazonaws.com/bucket-name/key-name
	// The MinIO client typically uses path-style for presigned URLs if not configured otherwise.
	// If `configuredBaseURL` is already a virtual-hosted style URL for a specific bucket, then `bucketName` should not be added.

	// Let's assume `configuredBaseURL` is a prefix that *may or may not* include the bucket.
	// If the `previewBaseURL` in `MinioS3Service` (set during construction) is intended to be the *actual*
	// base for object URLs (e.g. "https://mybucket.s3.amazonaws.com" or "https://mycdn.com/mybucket"),
	// then we just append objectKey.
	// If `s.previewBaseURL` is more generic (e.g. "https://s3.amazonaws.com"), then we append bucket and key.

	// MODIFIED: As per user request, bucketName is no longer part of the path construction here.
	// It's assumed that `configuredBaseURL` points to a location where appending `objectKey` is sufficient.
	// e.g., `configuredBaseURL` = "http://myservice.com/actual-bucket-path-if-needed"
	// or `configuredBaseURL` = "http://myservice.com" and `objectKey` = "bucket-simulating-path/item.jpg"
	base.Path = path.Join(base.Path, objectKey)
	return base.String(), nil

	// NOTE: For private objects, pre-signed URLs are the way to go.
	// Example for presigned GET URL (would require different logic and config):
	/*
	if s.client == nil {
		return "", fmt.Errorf("MinIO client is not initialized")
	}
	reqParams := make(url.Values)
	// reqParams.Set("response-content-disposition", "attachment; filename=\"your-filename.ext\"") // Optional
	presignedURL, err := s.client.PresignedGetObject(context.Background(), bucketName, objectKey, time.Second*24*60*60, reqParams) // 1 day expiry
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL for object '%s' in bucket '%s': %w", objectKey, bucketName, err)
	}
	return presignedURL.String(), nil
	*/
}
