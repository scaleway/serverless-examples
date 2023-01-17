package cors

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
)

// See: https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS#the_http_response_headers
func setCorsHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func HandleWithCors(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	fmt.Printf("%q", dump)

	// Sets the response headers to allow CORS requests.
	setCorsHeaders(w)

	w.Header().Set("Content-type", "text/plain")
	w.WriteHeader(http.StatusOK)

	_, err = io.WriteString(w, "This function is allowing most CORS requests")
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
	}
}
