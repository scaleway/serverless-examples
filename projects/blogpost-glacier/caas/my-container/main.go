package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Glacier is a struct to represent a Glacier entry, definied by an ID
// for database correlation and data to read the glacier image.
type Glacier struct {
	ID       int    // unique identifier of the ressource in database
	ImageURI string // way to access the image
}

// GlacierPageData is the struct that represents the displayed data containing
// metadata and Glaciers.
type GlacierPageData struct {
	PageTitle string    // title of the page to display
	Glaciers  []Glacier // list of glaciers to display
}

var (
	pgClient *sql.DB
	s3Client *minio.Client
)

func main() {
	initConnexions()

	// close database
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Printf("error during database closing : %v", err)
		}
	}(pgClient)

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/img", func(writer http.ResponseWriter, request *http.Request) {
		queryKeys, ok := request.URL.Query()["id"]

		if !ok || len(queryKeys[0]) == 0 {
			log.Println("Url Param 'id' is missing")
			writer.WriteHeader(http.StatusBadRequest)

			return
		}

		idInt, err := strconv.Atoi(queryKeys[0])
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusBadRequest)

			return
		}

		image, err := getFileFromS3(idInt)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)

			return
		}

		writer.Header().Set("Content-Type", "application/image")
		writer.Header().Set("Content-Disposition", "attachment; filename="+queryKeys[0]+".jpg")
		writer.WriteHeader(http.StatusOK)

		_, err = writer.Write(image)

		if err != nil {
			log.Fatal(err)
		}
	})

	// using go template is enough to display glacier data
	tmpl := template.Must(template.ParseFiles("layout.html"))

	glaciers, err := readGlaciers()
	if err != nil {
		log.Panicf("read glaciers error : %v", err)
	}

	http.HandleFunc("/", func(writer http.ResponseWriter, reader *http.Request) {
		data := GlacierPageData{
			PageTitle: "Glacier list " + time.Now().Format(time.RFC850),
			Glaciers:  glaciers,
		}

		if err := tmpl.Execute(writer, data); err != nil {
			log.Panicf("error template reading %v", err)
		}
	})

	listeningPort := os.Getenv("LISTENING_PORT")

	if listeningPort == "" {
		listeningPort = ":8080"
	}

	if err := http.ListenAndServe(listeningPort, nil); err != nil {
		log.Panicf("error listen and serve on port %s", listeningPort)
	}
}

// readGlaciers is used to get all the glaciers data from database.
func readGlaciers() ([]Glacier, error) {
	rows, err := pgClient.Query(`SELECT id, file_name FROM data`)
	if err != nil {
		log.Panicf("error query select database %v", err)
	}

	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			log.Panicf("error during rows closing : %v", err)
		}
	}(rows)

	glaciers := make([]Glacier, 0)

	for rows.Next() {
		var (
			id       int
			fileName string
		)

		if err := rows.Scan(&id, &fileName); err != nil {
			log.Panicf("error during rows scan : %v", err)
		}

		glaciers = append(glaciers, Glacier{
			ImageURI: fileName,
			ID:       id,
		})
	}

	return glaciers, nil
}

func getFileFromS3(id int) ([]byte, error) {
	// reuse readGlaciers to determine image name
	glaciers, err := readGlaciers()
	if err != nil {
		return nil, fmt.Errorf("error readGlaciers (db) %w", err)
	}

	imgName := ""

	for idx := range glaciers {
		if glaciers[idx].ID == id {
			imgName = glaciers[idx].ImageURI

			break
		}
	}

	reader, err := s3Client.GetObject(context.Background(), "tutoglacier", imgName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("error get object (image) s3 : %w", err)
	}
	defer reader.Close()

	localFile, err := os.CreateTemp("", "image-*.jpg")
	if err != nil {
		return nil, fmt.Errorf("error create temporary file  %w", err)
	}
	defer localFile.Close()

	stat, err := reader.Stat()
	if err != nil {
		return nil, fmt.Errorf("reader stat error %w", err)
	}

	if _, err := io.CopyN(localFile, reader, stat.Size); err != nil {
		return nil, fmt.Errorf("error copying file %w", err)
	}

	fileContent, err := os.ReadFile(localFile.Name())
	if err != nil {
		return nil, fmt.Errorf("error reading file content %w", err)
	}

	return fileContent, nil
}

func initConnexions() {
	// DB
	var (
		host     = os.Getenv("DB_HOST")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASS")
		dbname   = os.Getenv("DB_NAME")
	)

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Panicf("error reading db port %v", err)
	}

	const pgConnStr = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
	psqlconn := fmt.Sprintf(pgConnStr, host, port, user, password, dbname)

	// open database
	pgClient, err = sql.Open("postgres", psqlconn)
	if err != nil {
		log.Panicf("error opening databse : %v", err)
	}

	// check db
	if err = pgClient.Ping(); err != nil {
		log.Panicf("error database does not respond to ping : %v", err)
	}

	// S3
	var (
		endpoint        = os.Getenv("S3_ENDPOINT")
		accessKeyID     = os.Getenv("S3_ACCESSKEY")
		secretAccessKey = os.Getenv("S3_SECRET")
	)

	s3Client, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		log.Panicf("error opening s3 %v", err)
	}
}
