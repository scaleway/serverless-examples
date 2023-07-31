package cors

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const offlineTestingServer = "http://localhost:8080"

func TestHandleWithCors(t *testing.T) {

	resp, err := http.Get(offlineTestingServer)
	assert.Nil(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "text/plain", resp.Header.Get("Content-Type"))

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "This function is allowing most CORS requests", string(bodyBytes))

}
