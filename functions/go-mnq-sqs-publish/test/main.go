package main

import (
	scw "go-mnq-publish"

	"github.com/scaleway/serverless-functions-go/local"
)

func main() {
	local.ServeHandler(scw.Handle, local.WithPort(8080))
}
