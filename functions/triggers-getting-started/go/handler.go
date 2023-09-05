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

// This handle function comes frome our examples and is not modified at all.
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := factorial(n)
	fmt.Printf("go: factorial of %d is %s\n", n, result.String())

	_, err = io.WriteString(w, result.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
}
