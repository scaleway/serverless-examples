package myfunc

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleWithCors(t *testing.T) {

	offlineTestingServer := "http://localhost:8080"

	body := []byte(`{
		"username": "Jane Doe",
		"message": "Hello World"
	}`)
	resp, err := http.Post(offlineTestingServer, "application/json", bytes.NewBuffer(body))
	assert.Nil(t, err)

	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

}
