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

const (
	offlineTestingServer = "http://localhost:8080"
	fileToUpload         = "go.sum"
)

func TestHandle(t *testing.T) {

	form := new(bytes.Buffer)
	writer := multipart.NewWriter(form)
	fw, err := writer.CreateFormFile("data", filepath.Base(fileToUpload))
	assert.NoError(t, err)

	fd, err := os.Open(fileToUpload)
	assert.NoError(t, err)

	defer fd.Close()
	_, err = io.Copy(fw, fd)
	assert.NoError(t, err)

	writer.Close()

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, offlineTestingServer, form)
	assert.NoError(t, err)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	assert.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

}
