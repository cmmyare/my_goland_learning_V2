package routes

import (
	"github.com/cmmyare/restapi/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterMoviesRoutes(server *gin.Engine) {
	server.POST("/create_movies", controllers.CreateMovie)
	server.PUT("/update_movies/:id", controllers.UpdateMovie)
	server.PATCH("/par_update_movies/:id", controllers.PatchMovie)
	server.GET("/get_movies", controllers.GetMovies)
	server.DELETE("/delete_movie/:id",controllers.DeleteMovie)
	server.POST("/find_movie_by_name", controllers.FindMovieByName)
	server.POST("/find_movie_by_id", controllers.FindMovieByID)
}