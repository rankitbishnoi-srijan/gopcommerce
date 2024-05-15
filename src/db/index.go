package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect establishes a connection to the MongoDB database
func Connect(uri string) (*mongo.Client, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Disconnect closes the connection to the MongoDB database
func Disconnect(client *mongo.Client) error {
	err := client.Disconnect(context.Background())
	if err != nil {
		return err
	}

	return nil
}

// GetCollection returns a reference to a MongoDB collection
func GetCollection(client *mongo.Client, databaseName, collectionName string) *mongo.Collection {
	return client.Database(databaseName).Collection(collectionName)
}

// InsertOne inserts a single document into a MongoDB collection
func InsertOne(collection *mongo.Collection, document interface{}) (*mongo.InsertOneResult, error) {
	result, err := collection.InsertOne(context.Background(), document)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// FindOne retrieves a single document from a MongoDB collection
func FindOne(collection *mongo.Collection, filter interface{}) *mongo.SingleResult {
	return collection.FindOne(context.Background(), filter)
}

// Find retrieves multiple documents from a MongoDB collection
func Find(collection *mongo.Collection, filter interface{}) (*mongo.Cursor, error) {
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	return cursor, nil
}

// UpdateOne updates a single document in a MongoDB collection
func UpdateOne(collection *mongo.Collection, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteOne removes a single document from a MongoDB collection
func DeleteOne(collection *mongo.Collection, filter interface{}) (*mongo.DeleteResult, error) {
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteMany removes multiple documents from a MongoDB collection
func DeleteMany(collection *mongo.Collection, filter interface{}) (*mongo.DeleteResult, error) {
	result, err := collection.DeleteMany(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// CountDocuments returns the number of documents in a MongoDB collection
func CountDocuments(collection *mongo.Collection, filter interface{}) (int64, error) {
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Aggregate performs an aggregation operation on a MongoDB collection
func Aggregate(collection *mongo.Collection, pipeline interface{}) (*mongo.Cursor, error) {
	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}

	return cursor, nil
}
