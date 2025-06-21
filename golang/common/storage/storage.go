package storage

import "io"

// StorageService defines the interface for interacting with a storage backend.
// It's designed to be compatible with MinIO, Aliyun OSS (via S3 compatibility), and AWS S3.
type StorageService interface {
	// Upload uploads a file to the specified bucket with the given object key.
	// It takes an io.Reader for the file content, the object size, and content type.
	// It returns the path/key of the uploaded object or an error.
	Upload(bucketName string, objectKey string, reader io.Reader, objectSize int64, contentType string) (string, error)

	// Download retrieves a file from the specified bucket with the given object key.
	// It returns an io.ReadCloser for the file content or an error. The caller is responsible for closing the reader.
	Download(bucketName string, objectKey string) (io.ReadCloser, error)

	// GetPreviewURL generates a URL suitable for previewing an object.
	// For S3-compatible services, this might be a pre-signed URL or a direct public URL
	// if the object is public. The baseURL is the public-facing prefix for your storage
	// (e.g., "https://cdn.example.com" or "https://s3.region.amazonaws.com/bucketname").
	// If the objects are intended to be public, this method can simply concatenate baseURL, bucketName, and objectKey.
	// If pre-signed URLs are needed for private objects, this method would generate them.
	// For simplicity in this initial version, we'll assume public URLs or that the MinIO client
	// can provide a base for this. More complex pre-signing logic can be added if needed.
	GetPreviewURL(bucketName string, objectKey string, baseURL string) (string, error)
}
