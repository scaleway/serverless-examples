package handler

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandle(t *testing.T) {

	offlineTestingServer := "http://localhost:8080"
	fileToUpload := "go.sum"

	form := new(bytes.Buffer)
	writer := multipart.NewWriter(form)
	fw, err := writer.CreateFormFile("data", filepath.Base(fileToUpload))
	assert.Nil(t, err)

	fd, err := os.Open(fileToUpload)
	assert.Nil(t, err)

	defer fd.Close()
	_, err = io.Copy(fw, fd)
	assert.Nil(t, err)

	writer.Close()

	client := &http.Client{}
	req, err := http.NewRequest("POST", offlineTestingServer, form)
	assert.Nil(t, err)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	assert.Nil(t, err)

	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

}
