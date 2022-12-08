package myfunc

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Handle - Handle event
func Handle(w http.ResponseWriter, r *http.Request) {
	tmpImage, err := readImageToTempFile()
	if err != nil {
		log.Panicf("error on saving image : %v", err)
	}

	if err := connectAndPushImageToS3(tmpImage); err != nil {
		log.Panicf("error on uploading image to s3 %v", err)
	}

	if err := writeDataToDatabase(tmpImage); err != nil {
		log.Panicf("error on writing to database : %v", err)
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
}

// readImageToTempFile is used to read the image from the API and stores it in
// a temporary file (that can be deleted after usage).
// It returns a string (full path+name of the temp file) and an empty error if
// everything is ok.
func readImageToTempFile() (string, error) {
	// time now to generate URL to download webcam image
	timeNow := time.Now()

	const unformattedURL = "TODO" // SETUP YOUR URL HERE. Webcams often use a dynamic URL to retreive images in jpg.
	// for example we can have this kind of URL : https://webcam.dev/scaleway_mount/%d/%02d/%d.jpg

	// populate the URL with parameters
	urlDate := fmt.Sprintf(unformattedURL, timeNow.Year(), int32(timeNow.Month()), timeNow.Day()-1)

	// execute GET request to the URL (should return the image)
	response, err := http.Get(urlDate)
	if err != nil {
		return "", fmt.Errorf("error on retreiving url : %w", err)
	}
	defer response.Body.Close()

	// preapre temp image name
	const unformattedImageName = "temp-image-%d-%d-%d-*.jpg"

	imageName := fmt.Sprintf(unformattedImageName, timeNow.Year(), int32(timeNow.Month()), timeNow.Day())

	// create the temporary file to store image
	file, err := os.CreateTemp("", imageName)
	if err != nil {
		return "", fmt.Errorf("error on creating temp file %w", err)
	}
	defer file.Close()

	// copy the request body (image) in the temp file
	_, err = io.Copy(file, response.Body)

	return file.Name(), err
}

// connectAndPushImageToS3 will push "filePath" to the S3 bucket from env variables.
func connectAndPushImageToS3(filePath string) error {
	ctx := context.Background()

	// read vars from env
	var (
		endpoint        = os.Getenv("S3_ENDPOINT")
		accessKeyID     = os.Getenv("S3_ACCESSKEY")
		secretAccessKey = os.Getenv("S3_SECRET")
	)

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		return fmt.Errorf("error opening s3 connexion : %w", err)
	}

	// Make a new bucket called mymusic.
	const (
		bucketName = "tutoglacier"
		location   = "fr-par"
	)

	if err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location}); err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			return fmt.Errorf("error on checking bucket existance : %w", err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	const (
		contentType = "application/image"
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

// writeDataToDatabase will create a new record in the database to store data about
// tmpImage. In the future it can be used to add annotations to the images or for example
// add some context data such as glacier temperature (fetched via an API) etc...
func writeDataToDatabase(tmpImage string) error {
	// read vars from env
	var (
		host     = os.Getenv("DB_HOST")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASS")
		dbname   = os.Getenv("DB_NAME")
	)

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return fmt.Errorf("error on reading database port : %w", err)
	}

	// prepare connexion to postgres (we could use other databases)
	const pgConnStr = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
	psqlconn := fmt.Sprintf(pgConnStr, host, port, user, password, dbname)

	// open database
	dbGlacier, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return fmt.Errorf("error opening databse : %w", err)
	}

	// close database
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Printf("error during database closing : %v", err)
		}
	}(dbGlacier)

	// check db
	if err = dbGlacier.Ping(); err != nil {
		return fmt.Errorf("error database does not respond to ping : %w", err)
	}

	// create the table if it does not exists to avoid errors
	_, err = dbGlacier.Exec("CREATE TABLE IF NOT EXISTS data (id SERIAL, file_name TEXT NOT NULL, meta_data TEXT NULL)")
	if err != nil {
		return fmt.Errorf("error failed to create table if not exists  : %w", err)
	}

	// insert a new record to the database
	_, err = dbGlacier.Exec("INSERT INTO data (file_name) VALUES ($1)", filepath.Base(tmpImage))
	if err != nil {
		return fmt.Errorf("error failed to insert data : %w", err)
	}

	return nil
}
