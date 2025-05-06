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

// client is a global variable that holds the MongoDB client instance.
var client *mongo.Client

// RandomDocument represents a document in the MongoDB collection with an ID and a Name.
type RandomDocument struct {
	ID   int    `bson:"id"`
	Name string `bson:"name"`
}

// init initializes the MongoDB client connection.
// It reads the MongoDB connection details from environment variables and establishes a connection.
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

	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	slog.Info("Connected to MongoDB")
}

// Handle is an HTTP handler function that inserts a random document into the MongoDB collection
// and returns the inserted document as a JSON response.
func Handle(w http.ResponseWriter, _ *http.Request) {
	collection := client.Database("testdb").Collection("testcollection")

	randomID := rand.Intn(100000000)

	doc := RandomDocument{
		ID:   randomID,
		Name: fmt.Sprintf("RandomName%d", randomID),
	}

	insertResult, err := collection.InsertOne(context.TODO(), doc)
	if err != nil {
		http.Error(w, "Failed to insert document", http.StatusInternalServerError)
		return
	}

	slog.Info("Inserted document", "id", insertResult.InsertedID)

	var result RandomDocument

	if err := collection.FindOne(context.TODO(), bson.M{"_id": insertResult.InsertedID}).Decode(&result); err != nil {
		http.Error(w, "Failed to find document", http.StatusInternalServerError)
		return
	}

	slog.Info("Found document", "name", result.Name, "id", result.ID)

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
