package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/rankitbishnoi-srijan/gopcommerce/src/db"
	"github.com/rankitbishnoi-srijan/gopcommerce/src/utils"
)

func NewRouter(db *mongo.Client) *gin.Engine {
	// Set the router as the default one shipped with Gin
	router := gin.New()
	expectedHost := "localhost:8080"

	// Setup middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(ApiMiddleware(db))

	// Setup Security Headers
	router.Use(func(c *gin.Context) {
		if c.Request.Host != expectedHost {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid host header"})
			return
		}
		c.Header("X-Frame-Options", "DENY")
		c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		c.Header("Referrer-Policy", "strict-origin")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")
		c.Next()
	})

	// Serve frontend static files
	router.Use(static.Serve("/", static.LocalFile("./public/", true)))

	// Setup route group for the API
	api := router.Group("/api")
	{
		apiHandler := func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Uniform API",
			})
		}
		api.GET("", apiHandler)
		api.GET("/", apiHandler)
	}

	return router
}

// ApiMiddleware will add the db connection to the context
func ApiMiddleware(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("databaseConn", db)
		c.Next()
	}
}

func main() {
	// Load environment variables
	if err := utils.LoadEnv(); err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	// Initialize MongoDB
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI must be set")
	}
	client, err := db.Connect(uri)

	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	defer db.Disconnect(client)

	// Initialize router
	router := NewRouter(client)

	// Create server with timeout
	srv := &http.Server{
		Addr:    ":" + httpPort,
		Handler: router,
		// set timeout due CWE-400 - Potential Slowloris Attack
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Printf("Failed to start server: %v", err)
	}
}
