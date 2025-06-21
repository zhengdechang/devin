package storage

import (
	"fmt"
)

// Config holds the necessary configuration parameters to initialize a StorageService.
// It's designed to be flexible enough for MinIO, AWS S3, and Aliyun OSS (via S3 compatibility).
type Config struct {
	// Type specifies the backend type, though with a unified MinioS3Service, this might be less critical
	// unless other types are added later. For now, it's implicitly "s3compatible".
	// Type string `json:"type"` // Example: "minio", "s3", "aliyun_oss_s3_compatible"

	Endpoint        string `json:"endpoint"`         // e.g., "s3.amazonaws.com", "oss-cn-hangzhou.aliyuncs.com", "localhost:9000"
	AccessKeyID     string `json:"accessKeyID"`      // Access Key ID
	SecretAccessKey string `json:"secretAccessKey"`  // Secret Access Key
	UseSSL          bool   `json:"useSSL"`           // True to use HTTPS, false for HTTP
	Region          string `json:"region"`           // AWS region, e.g., "us-east-1". Important for S3. Optional for MinIO.
	DefaultBucket   string `json:"defaultBucket"`    // Optional: A default bucket to use if not specified in API calls. (Not directly used by MinioS3Service methods yet)

	// BaseURL for constructing preview URLs. This is the public-facing prefix for your storage.
	// Examples:
	// - For objects served via a CDN: "https://my-cdn.com"
	// - For direct S3 (path-style, bucket in path): "https://s3.region.amazonaws.com"
	// - For direct S3 (virtual-hosted, bucket in hostname): "https://my-bucket.s3.region.amazonaws.com"
	// - For Aliyun OSS (path-style): "https://oss-cn-hangzhou.aliyuncs.com"
	// - For Aliyun OSS (virtual-hosted): "https://my-bucket.oss-cn-hangzhou.aliyuncs.com"
	// - For MinIO: "http://localhost:9000" (if API endpoint is also base for preview) or "http://public-minio-ip:9000"
	// The GetPreviewURL method will append bucketName (if not part of BaseURL) and objectKey to this.
	PreviewBaseURL  string `json:"previewBaseURL"`
}

// NewStorageService is a factory function that creates and returns a StorageService
// based on the provided configuration. Currently, it only supports the MinioS3Service.
func NewStorageService(cfg Config) (StorageService, error) {
	if cfg.Endpoint == "" {
		return nil, fmt.Errorf("storage endpoint must be configured")
	}
	if cfg.AccessKeyID == "" {
		return nil, fmt.Errorf("storage accessKeyID must be configured")
	}
	if cfg.SecretAccessKey == "" {
		return nil, fmt.Errorf("storage secretAccessKey must be configured")
	}
	// PreviewBaseURL is critical for GetPreviewURL as currently implemented.
	// If allowing pre-signed URLs in the future, this might become optional.
	if cfg.PreviewBaseURL == "" {
		// Depending on requirements, this could be a warning or an error.
		// If GetPreviewURL is not always used, or if an alternative method for URL generation exists (like pre-signed),
		// this check might be relaxed. For the current GetPreviewURL, it's essential.
		return nil, fmt.Errorf("storage previewBaseURL must be configured for GetPreviewURL functionality")
	}


	// For now, we only have one type of service, MinioS3Service, which handles all S3-compatible backends.
	// The MinioS3Service constructor takes region and previewBaseURL directly.
	service, err := NewMinioS3Service(
		cfg.Endpoint,
		cfg.AccessKeyID,
		cfg.SecretAccessKey,
		cfg.Region,         // Pass region
		cfg.PreviewBaseURL, // Pass previewBaseURL
		cfg.UseSSL,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO S3 compatible service: %w", err)
	}

	return service, nil
}
