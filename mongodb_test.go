package messagesStorageMongo

import (
	"context"
	"github.com/mailhedgehog/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func createMongoDbConnection() *mongo.Database {
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017").SetTimeout(5 * time.Second)

	clientOptions = clientOptions.SetAuth(options.Credential{
		Username: "test_root",
		Password: "test_secret",
	})

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	logger.PanicIfError(err)

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	logger.PanicIfError(err)

	logManager().Debug("Connected to MongoDB")

	return client.Database("test_db")
}

func createMongoTestCollection() *mongo.Collection {
	collection := createMongoDbConnection().Collection("bar")

	// Truncate in case data saved form previous test.
	collection.DeleteMany(context.TODO(), bson.D{})

	return collection
}
