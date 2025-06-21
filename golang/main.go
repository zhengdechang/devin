package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"storage-module/common/storage" // Adjust import path based on your module name in go.mod
)

func main() {
	fmt.Println("Storage Service Module Example")
	fmt.Println("==============================")

	// This example demonstrates how to configure the service for different backends.
	// IMPORTANT: Replace placeholder values with your actual credentials and endpoints.
	// For security, these should come from environment variables or a config file, not hardcoded.

	// Common setup for the example
	exampleBucketName := "my-unique-test-bucket-jules" // Change to a unique bucket name
	exampleObjectKey := "test-files/my-sample-object.txt"
	exampleContent := "Hello, S3-Compatible World! This is a test content."

	// --- Configuration for MinIO (Self-hosted or Playground) ---
	minioConfig := storage.Config{
		Endpoint:        "localhost:9000", // Or your MinIO server endpoint like "play.min.io"
		AccessKeyID:     "YOUR_MINIO_ACCESS_KEY",    // e.g., "minioadmin" for local default
		SecretAccessKey: "YOUR_MINIO_SECRET_KEY", // e.g., "minioadmin" for local default
		UseSSL:          false,                   // Set to true if your MinIO uses SSL
		Region:          "",                      // Usually not required for MinIO unless specific setup
		// PreviewBaseURL should be the base HTTP address where objects can be accessed.
		// GetPreviewURL will now append the objectKey directly to this URL.
		// Example: If PreviewBaseURL = "http://localhost:9000/my-bucket-alias",
		// and objectKey = "public/file.txt",
		// result = "http://localhost:9000/my-bucket-alias/public/file.txt".
		// If MinIO serves 'exampleBucketName' directly at "http://localhost:9000/exampleBucketName/",
		// and you want URLs like "http://localhost:9000/exampleBucketName/public/file.txt",
		// then set PreviewBaseURL = "http://localhost:9000/exampleBucketName".
		PreviewBaseURL:  fmt.Sprintf("http://localhost:9000/%s", exampleBucketName), // Assumes bucket is part of the base path
	}

	// --- Configuration for AWS S3 ---
	awsS3Config := storage.Config{
		Endpoint:        "s3.YOUR_AWS_REGION.amazonaws.com", // e.g., "s3.us-east-1.amazonaws.com"
		AccessKeyID:     "YOUR_AWS_ACCESS_KEY_ID",
		SecretAccessKey: "YOUR_AWS_SECRET_ACCESS_KEY",
		UseSSL:          true,
		Region:          "YOUR_AWS_REGION", // e.g., "us-east-1"
		// PreviewBaseURL for S3:
		// - Virtual-hosted style (recommended for new buckets): "https://YOUR_BUCKET_NAME.s3.YOUR_AWS_REGION.amazonaws.com"
		// - Path-style (older, may be deprecated for new regions/buckets): "https://s3.YOUR_AWS_REGION.amazonaws.com/YOUR_BUCKET_NAME"
		// GetPreviewURL will append the objectKey to this.
		PreviewBaseURL:  fmt.Sprintf("https://%s.s3.%s.amazonaws.com", exampleBucketName, "YOUR_AWS_REGION"), // Virtual-hosted style base
		// Or for path-style:
		// PreviewBaseURL:  fmt.Sprintf("https://s3.%s.amazonaws.com/%s", "YOUR_AWS_REGION", exampleBucketName),
	}


	// --- Configuration for Aliyun OSS (S3 Compatible Mode) ---
	aliyunOSSConfig := storage.Config{
		Endpoint:        "oss-YOUR_OSS_REGION.aliyuncs.com", // e.g., "oss-cn-hangzhou.aliyuncs.com"
		AccessKeyID:     "YOUR_ALIYUN_ACCESS_KEY_ID",
		SecretAccessKey: "YOUR_ALIYUN_SECRET_ACCESS_KEY",
		UseSSL:          true,
		Region:          "YOUR_OSS_REGION", // e.g., "cn-hangzhou" (often part of endpoint too)
		// PreviewBaseURL for Aliyun OSS:
		// - Virtual-hosted style: "https://YOUR_BUCKET_NAME.oss-YOUR_OSS_REGION.aliyuncs.com"
		// - Path-style: "https://oss-YOUR_OSS_REGION.aliyuncs.com/YOUR_BUCKET_NAME"
		// GetPreviewURL will append the objectKey to this.
		PreviewBaseURL:  fmt.Sprintf("https://%s.oss-%s.aliyuncs.com", exampleBucketName, "YOUR_OSS_REGION"), // Virtual-hosted style base
		// Or for path-style:
		// PreviewBaseURL: fmt.Sprintf("https://oss-%s.aliyuncs.com/%s", "YOUR_OSS_REGION", exampleBucketName),
	}

	// Select the configuration to use for testing.
	// !!! IMPORTANT: Fill in the credentials for the config you want to test.
	// By default, this example will not run unless you update placeholder credentials.
	// For this example, we'll try to run with a dummy "minioConfig"
	// but expect it to fail if not configured.

	// currentConfig := minioConfig
	// currentConfigName := "MinIO"

	// currentConfig := awsS3Config
	// currentConfigName := "AWS S3"

	 currentConfig := aliyunOSSConfig
	 currentConfigName := "Aliyun OSS (S3 Compatible)"


	// Check if placeholder values are still there for the selected config
	if strings.HasPrefix(currentConfig.AccessKeyID, "YOUR_") || currentConfig.Endpoint == "" || strings.Contains(currentConfig.Endpoint, "YOUR_") {
		log.Printf("Placeholder credentials or endpoint found for %s. Skipping actual test execution.", currentConfigName)
		log.Println("Please update main.go with your actual credentials and endpoint to run the example.")
		log.Printf("Selected config: %+v\n", currentConfig)
		os.Exit(0) // Exit gracefully, as we can't run the test.
	}

	log.Printf("Attempting to use configuration for: %s\n", currentConfigName)
	log.Printf("Endpoint: %s, Bucket: %s, ObjectKey: %s\n", currentConfig.Endpoint, exampleBucketName, exampleObjectKey)


	// Create a new storage service instance
	svc, err := storage.NewStorageService(currentConfig)
	if err != nil {
		log.Fatalf("Failed to create storage service for %s: %v", currentConfigName, err)
	}

	log.Printf("Successfully created storage service for %s.", currentConfigName)

	// --- Test Upload ---
	log.Printf("Attempting to upload '%s' to bucket '%s'...", exampleObjectKey, exampleBucketName)
	contentReader := bytes.NewReader([]byte(exampleContent))
	uploadedKey, err := svc.Upload(exampleBucketName, exampleObjectKey, contentReader, int64(len(exampleContent)), "text/plain")
	if err != nil {
		log.Fatalf("Failed to upload to %s: %v", currentConfigName, err)
	}
	log.Printf("Successfully uploaded object. Returned key: %s\n", uploadedKey)

	// --- Test GetPreviewURL ---
	log.Printf("Attempting to get preview URL for '%s' in bucket '%s'...", exampleObjectKey, exampleBucketName)
	// The GetPreviewURL in MinioS3Service uses the PreviewBaseURL from its own configuration (which was set from currentConfig.PreviewBaseURL)
	// So, the third argument to GetPreviewURL in the interface (baseURL) is effectively s.previewBaseURL.
	previewURL, err := svc.GetPreviewURL(exampleBucketName, exampleObjectKey, currentConfig.PreviewBaseURL)
	if err != nil {
		log.Fatalf("Failed to get preview URL from %s: %v", currentConfigName, err)
	}
	log.Printf("Successfully generated preview URL: %s\n", previewURL)
	log.Println("Note: This URL's accessibility depends on bucket/object permissions and correct PreviewBaseURL configuration.")

	// --- Test Download ---
	log.Printf("Attempting to download '%s' from bucket '%s'...", exampleObjectKey, exampleBucketName)
	downloadedObject, err := svc.Download(exampleBucketName, exampleObjectKey)
	if err != nil {
		log.Fatalf("Failed to download from %s: %v", currentConfigName, err)
	}
	defer downloadedObject.Close()

	downloadedContent, err := io.ReadAll(downloadedObject)
	if err != nil {
		log.Fatalf("Failed to read downloaded content from %s: %v", currentConfigName, err)
	}

	if string(downloadedContent) == exampleContent {
		log.Printf("Successfully downloaded and verified content: %s\n", string(downloadedContent))
	} else {
		log.Fatalf("Content mismatch! Expected: '%s', Got: '%s'", exampleContent, string(downloadedContent))
	}

	log.Printf("\nExample operations for %s completed successfully!\n", currentConfigName)
	log.Println("Remember to clean up any created buckets or objects in your storage provider if necessary.")


	// --- How to run this example ---
	// 1. Ensure you have Go installed.
	// 2. Navigate to the `golang` directory in your terminal: `cd golang`
	// 3. Update the placeholder credentials in `main.go` for the service you want to test.
	// 4. Run `go mod tidy` to ensure dependencies are correctly managed.
	// 5. Run the example: `go run main.go`
	//
	// Example for creating a local dummy file to upload (optional, if you change Upload to use a file path)
	/*
	 tempFilePath := filepath.Join(os.TempDir(), "upload_test_file.txt")
	 if err := os.WriteFile(tempFilePath, []byte("This is a test file for upload."), 0644); err != nil {
		log.Fatalf("Failed to create temp file: %v", err)
	 }
	 log.Printf("Created temp file for upload: %s", tempFilePath)
	 // To upload from file:
	 // file, err := os.Open(tempFilePath)
	 // if err != nil { log.Fatalf("Failed to open file: %v", err) }
	 // defer file.Close()
	 // fileInfo, _ := file.Stat()
	 // uploadedKey, err = svc.Upload(exampleBucketName, exampleObjectKey, file, fileInfo.Size(), "text/plain")
	 // ...
	 // os.Remove(tempFilePath) // Clean up
	*/
}
