package main

import (
	scw "github.com/scaleway/serverless-examples/functions/go-mongo"
	"github.com/scaleway/serverless-functions-go/local"
)

// https://github.com/scaleway/serverless-functions-go
func main() {
	// Replace "Handle" with your function handler name if necessary
	local.ServeHandler(scw.Handle, local.WithPort(8080))
}
