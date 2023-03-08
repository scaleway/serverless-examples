package handler

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "multipart/from-data" {
		http.Error(w, "this func only support multipart/from-data type", http.StatusBadRequest)

		return
	}
	// Entire files are sent as "multipart/form-data" so we need first to parse
	// the request before accessing file content.
	if err := r.ParseMultipartForm(32<<20 + 1024); err != nil {
		panic(err)
	}

	// Save the multipartForm of the request.
	multipartFile := r.MultipartForm.File["data"][0]

	// Open the file.
	file, err := multipartFile.Open()
	if err != nil {
		panic(err)
	}

	// Here we read the whole file content into fileContent as byte slice,
	// We can directly use the file descriptor to send it to s3 to avoid saving the file locally.
	fileContent, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	// Reach the temp dir of system to save the file
	tmpDir := os.TempDir()

	// Build the temp file path
	tempFile := filepath.Join(tmpDir, multipartFile.Filename)

	// Write file to system
	if err := os.WriteFile(tempFile, fileContent, 0o644); err != nil {
		panic(err)
	}

	// Always remove the file
	defer func() {
		os.Remove(tempFile)
	}()

	// If S3 enabled, send the file
	if strings.ToLower(os.Getenv("S3_ENABLED")) == "true" {
		log.Println("S3 upload enabled")

		if err := connectToS3AndPushFile(r.Context(), tempFile); err != nil {
			panic(err)
		}
	} else {
		log.Println("S3 upload disabled")
	}
}

// connectToS3AndPushFile will push file in parameter to the S3 bucket. Env variables (secrets)
// must be initialized to be able to upload files.
func connectToS3AndPushFile(ctx context.Context, filePath string) error {
	// Read required s3 variables
	var (
		endpoint        = os.Getenv("S3_ENDPOINT")
		accessKeyID     = os.Getenv("S3_ACCESSKEY")
		secretAccessKey = os.Getenv("S3_SECRET")
		bucketName      = os.Getenv("S3_BUCKET_NAME")
		region          = os.Getenv("S3_REGION")
	)

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		return fmt.Errorf("error opening s3 connexion : %w", err)
	}

	if err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: region}); err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("%s already exists \n", bucketName)
		} else {
			return fmt.Errorf("error on checking bucket status : %w", err)
		}
	} else {
		log.Printf("successfully created %s\n", bucketName)
	}

	const (
		contentType = "application/octet-stream"
	)

	info, err := minioClient.FPutObject(
		ctx,
		bucketName,
		filepath.Base(filePath),
		filePath,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return fmt.Errorf("error on put object to S3 bucket : %w", err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", filePath, info.Size)

	return nil
}
