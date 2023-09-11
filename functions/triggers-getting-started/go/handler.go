package handler

import (
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
)

func factorial(n int) *big.Int {
	var f big.Int
	f.MulRange(1, int64(n))
	return &f
}

func Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// SQS triggers are sent as POST requests.
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// The SQS trigger sends the message content in the body.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	n, err := strconv.Atoi(string(body))
	if err != nil {
		// Setting the status code to 200 will mark the message as processed.
		http.Error(w, err.Error(), http.StatusOK)
		return
	}

	result := factorial(n)
	fmt.Printf("go: factorial of %d is %s\n", n, result.String())

	// Because triggers are asynchronous, the response body is ignored.
	// It's kept here when testing locally.
	_, err = io.WriteString(w, result.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If the status code is not in the 2XX range, the message is considered
	// failed and is retried. In total, there are 3 retries.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
}
