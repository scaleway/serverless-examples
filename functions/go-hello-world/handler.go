package hello

import (
	"fmt"
	"io"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "Hello from a Go function!")

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
	}
}
