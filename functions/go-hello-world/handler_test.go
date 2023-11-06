package cors

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const offlineTestingServer = "http://localhost:8080"

func TestHelloWorld(t *testing.T) {
	resp, err := http.Get(offlineTestingServer)
	assert.NoError(t, err)

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "Hello from a Go function!", string(bodyBytes))
}
