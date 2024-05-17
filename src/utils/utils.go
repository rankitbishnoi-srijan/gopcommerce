package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"

	"crypto/md5"
	"encoding/hex"
)

// loadEnv loads environment variables from a .env file
func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	return nil
}

// GetDbClient returns the MongoDB client from the context
func GetDbClient(c *gin.Context) *mongo.Client {
	client := c.MustGet("databaseConn").(*mongo.Client)
	return client
}

// HashPassword hashes a password md5
func HashPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}
