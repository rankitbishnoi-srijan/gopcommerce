package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rankitbishnoi-srijan/gopcommerce/src/packages/auth"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/rankitbishnoi-srijan/gopcommerce/src/packages/users/models"
)

func getUsers(c *gin.Context) {

	collection := models.GetUserCollection(c)

	// Find all users
	cursor, err := collection.Find(c, bson.M{})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(c)

	// Iterate over the cursor
	var users []models.UserResponse
	for cursor.Next(c) {
		var user models.User
		err := cursor.Decode(&user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		users = append(users, models.NewUserResponse(user))
	}

	// Return the users
	c.JSON(http.StatusOK, models.UserListResponse{Users: users})
}

func createUser(c *gin.Context) {
	collection := models.GetUserCollection(c)

	// Bind the JSON data
	var user models.User
	err := c.BindJSON(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert the user
	_, err = collection.InsertOne(c, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the user
	c.JSON(http.StatusCreated, models.NewUserResponse(user))
}

func getUser(c *gin.Context) {
	collection := models.GetUserCollection(c)

	// Get the user ID
	id := c.Param("id")

	// Find the user
	filter := bson.M{"_id": id}
	user := collection.FindOne(c, filter)
	if user.Err() != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Return the user
	var userData models.User
	user.Decode(&userData)
	c.JSON(http.StatusOK, models.NewUserResponse(userData))
}

func deleteUser(c *gin.Context) {
	collection := models.GetUserCollection(c)

	// Get the user ID
	id := c.Param("id")

	// Find the user
	filter := bson.M{"_id": id}
	user := collection.FindOne(c, filter)
	if user.Err() != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Delete the user
	_, err := collection.DeleteOne(c, filter)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the user
	var userData models.User
	user.Decode(&userData)
	c.JSON(http.StatusOK, models.NewUserResponse(userData))
}

// RegisterHandlers registers the user handlers
func RegisterHandlers(router *gin.RouterGroup) {
	router.GET("/users", auth.AuthMiddleware(true), getUsers)
	router.POST("/users", auth.AuthMiddleware(true), createUser)
	router.GET("/users/:id", auth.AuthMiddleware(true), getUser)
	router.DELETE("/users/:id", auth.AuthMiddleware(true), deleteUser)
}
