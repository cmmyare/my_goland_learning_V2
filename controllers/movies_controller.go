package controllers

import (
	"net/http"
	"strings"

	"github.com/cmmyare/restapi/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateMovie(context *gin.Context) {
	var movie models.Movie
	if err := context.ShouldBindJSON(&movie); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := models.InserMovie(movie)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create movie"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Movie created successfully"})
}
// func UpdateMovie(context *gin.Context) {
// 	movieID := context.Param("id")
// 	var movie models.Movie
// 	if err := context.ShouldBindJSON(&movie); err != nil {
// 		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	err := models.UpdateMovie(movieID, movie)
// 	if err != nil {
// 		if strings.Contains(err.Error(), "no document found") {
// 			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
// 			return
// 		}
// 		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie"})
// 		return
// 	}
// 	context.JSON(http.StatusOK, gin.H{"message": "Movie updated successfully"})
// }

func UpdateMovie(context *gin.Context) {
	movieID := context.Param("id")
	var updateFields map[string]interface{}
	if err := context.ShouldBindJSON(&updateFields); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := models.UpdateMovie(movieID, updateFields)
	if err != nil {
		if strings.Contains(err.Error(), "no document found") {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Movie updated successfully"})
}


func GetMovies(Context *gin.Context) {
	moveies, err := models.ListAll()
	if err != nil {
		Context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
		return
	}
	if len(moveies) == 0 {
		Context.JSON(http.StatusNotFound, gin.H{"message": "no movies found"})
		return
	}
	Context.JSON(http.StatusOK, moveies)

}
// find movie by name
func FindMovieByName(context *gin.Context) {
    var req struct {
        MovieName string `json:"movie"`
    }
    if err := context.ShouldBindJSON(&req); err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }
    movie, err := models.FindByName(req.MovieName)

     if err != nil {
        if err == mongo.ErrNoDocuments {
            context.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
        } else {
            context.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        }
        return
    }
    context.JSON(http.StatusOK, movie)
}
func FindMovieByID(context *gin.Context) {
	var req struct {
		ID string `json:"_id"` // Expecting ObjectId as string in body
	}

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	movie, err := models.FindByID(req.ID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			context.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		} else if strings.Contains(err.Error(), "invalid ObjectID") {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		} else {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	context.JSON(http.StatusOK, movie)
}


func DeleteMovie(context *gin.Context) {
	movieID := context.Param("id")
	err := models.DeleteMovie(movieID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Movie deleted successfully"})
}

func PatchMovie(context *gin.Context) {
	movieID := context.Param("id")
	var payload map[string]interface{}
	if err := context.ShouldBindJSON(&payload); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Allow only known keys
	allowed := map[string]bool{"movie": true, "actors": true}
	filtered := make(map[string]interface{})
	for k, v := range payload {
		if allowed[k] {
			filtered[k] = v
		}
	}

	if len(filtered) == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "no valid fields provided (allowed: movie, actors)"})
		return
	}

	if err := models.PartialUpdateMovie(movieID, filtered); err != nil {
		if strings.Contains(err.Error(), "no document found") {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Movie partially updated successfully"})
}