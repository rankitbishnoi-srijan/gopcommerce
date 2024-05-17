package models

import (
	"github.com/gin-gonic/gin"
	"github.com/rankitbishnoi-srijan/gopcommerce/src/constants"
	"github.com/rankitbishnoi-srijan/gopcommerce/src/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// User represents a user in the system
type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username"`
	Password string             `json:"password" bson:"password"`
	Email    string             `json:"email" bson:"email"`
	Admin    bool               `json:"admin" bson:"admin"`
}

// UserLogin represents a user login request
type UserLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserRegister represents a user registration request
type UserRegister struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Admin    bool   `json:"admin"`
}

// UserResponse represents a user response
type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Admin    bool   `json:"admin"`
}

// NewUserResponse creates a new user response
func NewUserResponse(user User) UserResponse {
	return UserResponse{
		ID:       user.ID.Hex(),
		Username: user.Username,
		Email:    user.Email,
		Admin:    user.Admin,
	}
}

// UserListResponse represents a list of users response
type UserListResponse struct {
	Users []UserResponse `json:"users"`
}

// user collection
func GetUserCollection(c *gin.Context) *mongo.Collection {
	client := utils.GetDbClient(c)

	// Check if the user exists
	collection := client.Database(constants.DATABASE).Collection(constants.USERS_COLLECTION)

	return collection
}
