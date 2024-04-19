package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/scaleway/serverless-examples/containers/grpc-http2-go"
	"google.golang.org/grpc"
)

const containerEndpoint = "YOUR_CONTAINER_ENDPOINT:80"

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(containerEndpoint, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "Scaleway"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %s", r.GetMessage())
}
