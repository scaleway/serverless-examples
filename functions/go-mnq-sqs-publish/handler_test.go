package myfunc

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const offlineTestingServer = "http://localhost:8080"

func TestHandle(t *testing.T) {

	body := []byte(`{
		"username": "Jane Doe",
		"message": "Hello World"
	}`)
	resp, err := http.Post(offlineTestingServer, "application/json", bytes.NewBuffer(body))
	assert.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

}
