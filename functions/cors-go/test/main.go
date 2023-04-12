package main

import (
	handler "cors-go"

	"github.com/scaleway/serverless-functions-go/local"
)

func main() {
	// Replace "Handle" with your function handler name if necessary
	local.ServeHandler(handler.HandleWithCors, local.WithPort(8080))
}
