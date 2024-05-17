package auth

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/rankitbishnoi-srijan/gopcommerce/src/packages/users/models"
	"github.com/rankitbishnoi-srijan/gopcommerce/src/utils"
)

// AuthMiddleware is a middleware that checks if the request is authorized
func AuthMiddleware(checkAdmin ...bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}

			// Return the secret key
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Check if the token is valid
		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Get the user ID from the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		userID, ok := claims["user_id"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		collection := models.GetUserCollection(c)
		filter := map[string]interface{}{"_id": userID}
		user := collection.FindOne(c, filter)
		if user.Err() != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Check if the user is an admin
		if len(checkAdmin) > 0 && checkAdmin[0] {
			var userData models.User
			user.Decode(&userData)
			if !userData.Admin {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				return
			}
		}

		// Set the user ID in the context
		c.Set("user_id", userID)

		// Continue
		c.Next()
	}
}

// LoginHandler is a handler that authenticates the user
func LoginHandler(c *gin.Context) {
	// Get the username and password from the request
	username := c.PostForm("username")
	password := c.PostForm("password")

	collection := models.GetUserCollection(c)
	filter := map[string]interface{}{"username": username, "password": utils.HashPassword(password)}
	userData := collection.FindOne(c, filter)
	if userData.Err() != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get the user ID
	var user models.User
	userData.Decode(&user)
	userID := user.ID.Hex()

	// Create a new JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Fatalf("Failed to sign token: %v", err)
	}

	// Return the token
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// LogoutHandler is a handler that logs out the user
func LogoutHandler(c *gin.Context) {
	// Continue
	c.Next()
}

// RegisterHandlers registers the auth handlers
func RegisterHandlers(router *gin.RouterGroup) {
	router.POST("/login", LoginHandler)
	router.POST("/logout", AuthMiddleware(), LogoutHandler)
}
