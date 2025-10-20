// controllers/task_controller.go
package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cmmyare/restapi/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateTask(context *gin.Context) {
	var req models.Task
	// Bind JSON body
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}
	fmt.Println("req after Json >>>>>>>>>>>>>>>>>>>>>>>", req)
	fmt.Println("req.BookingInfo >>>>>>>>>>>>>>>>>>>>>>>", req.BookingInfo)

	// --- Convert string IDs to ObjectIDs ---
	// userID, err := primitive.ObjectIDFromHex(req.UserID.Hex())
	// fmt.Println("req.UserID without Hex() >>>>>>>>>>>>>>>>>>>>>>>", req.UserID)
	// if err != nil {
	// 	context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
	// 	return
	// }

	// taskerID, err := primitive.ObjectIDFromHex(req.TaskerID.Hex())
	// if err != nil {
	// 	context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tasker ID"})
	// 	return
	// }

	// categoryID, err := primitive.ObjectIDFromHex(req.CategoryID.Hex())
	// if err != nil {
	// 	context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
	// 	return
	// }

	// serviceLevelID, err := primitive.ObjectIDFromHex(req.ServiceLevelID.Hex())
	// if err != nil {
	// 	context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service level ID"})
	// 	return
	// }

	// --- Create Task struct with converted IDs ---
	task := models.Task{
		ID:             primitive.NewObjectID(),
		UserID:         req.UserID,
		TaskerID:       req.TaskerID,
		CategoryID:     req.CategoryID,
		ServiceLevelID: req.ServiceLevelID,
		BookingInfo:    req.BookingInfo,
	}

	// Optional: default timestamps or validations
	if task.BookingInfo.BookingDate.IsZero() {
		task.BookingInfo.BookingDate = time.Now()
	}

	// --- Insert into MongoDB ---
	if err := models.InsertTask(task); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Task created successfully",
		"task_id": task.ID.Hex(),
	})
}
