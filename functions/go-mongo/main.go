package example

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type RandomDocument struct {
	ID   int    `bson:"id"`
	Name string `bson:"name"`
}

func init() {
	mongoPublicEndpoint := os.Getenv("MONGO_PUBLIC_ENDPOINT")
	if mongoPublicEndpoint == "" {
		panic("MONGO_PUBLIC_ENDPOINT is required")
	}

	mongoUser := os.Getenv("MONGO_USER")
	if mongoUser == "" {
		panic("MONGO_USER is required")
	}

	mongoPassword := url.PathEscape(os.Getenv("MONGO_PASSWORD"))
	if mongoPassword == "" {
		panic("MONGO_PASSWORD is required")
	}

	// This is a basic sample that does not use certificate for authentication, not recommended for production.
	mongoURI := fmt.Sprintf(`mongodb+srv://%s:%s@%s/?tls=true&tlsInsecure=true`,
		mongoUser,
		mongoPassword,
		mongoPublicEndpoint)

	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongo_uri))
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	slog.Info("Connected to MongoDB")
}

func Handle(w http.ResponseWriter, r *http.Request) {
	collection := client.Database("testdb").Collection("testcollection")

	randomID := rand.Intn(100000000)

	doc := RandomDocument{
		ID:   randomID,
		Name: fmt.Sprintf("RandomName%d", randomID),
	}

	insertResult, err := collection.InsertOne(context.TODO(), doc)
	if err != nil {
		panic(err)
	}

	slog.Info("Inserted document with ID:", insertResult.InsertedID)

	var result RandomDocument

	if err := collection.FindOne(context.TODO(), bson.M{"_id": insertResult.InsertedID}).Decode(&result); err != nil {
		panic(err)
	}

	slog.Info("Found document", result.Name, result.ID)

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
}
