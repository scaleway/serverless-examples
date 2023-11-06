package main

import (
	handler "hello-go"

	"github.com/scaleway/serverless-functions-go/local"
)

func main() {
	local.ServeHandler(handler.Handle, local.WithPort(8080))
}
