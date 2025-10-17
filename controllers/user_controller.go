package controllers

import (
	"net/http"

	"github.com/cmmyare/restapi/models"
	"github.com/cmmyare/restapi/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)
type CreateUserResponse struct {
	Message  string `json:"message"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func CreateUser(context *gin.Context) {
	var user models.User

	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}

	if user.Username == "" || user.Email == "" || user.Password == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Username, Email, and Password are required"})
		return
	}

	collection := models.MongoClient.Database(models.DB).Collection("users")
	filter := bson.M{"email": user.Email}

	count, err := collection.CountDocuments(context, filter)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Database error while checking email"})
		return
	}
	if count > 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = hashedPassword

	err = models.InsertUser(user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

context.JSON(http.StatusOK, CreateUserResponse{
	Message:  "User created successfully",
	Username: user.Username,
	Email:    user.Email,
})
}

func LoginUser(context *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse request body
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Find user by email
	user, err := models.FindUserByEmail(req.Email)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare passwords
	if err := utils.CompareHashAndPassword(user.Password, req.Password); err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Success
	context.JSON(http.StatusOK, CreateUserResponse{
	Message:  "Login successful",
	Username: user.Username,
	Email:    user.Email,
})
}
