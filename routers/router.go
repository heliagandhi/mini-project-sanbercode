package routers

import (
	"mini-project-sanbercode/controllers"

	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	router := gin.Default()

	router.POST("/bioskop", controllers.CreateBioskop)
	router.GET("/bioskop", controllers.GetBioskops)
	router.GET("/bioskop/:id", controllers.GetBioskopByID)
	router.PUT("/bioskop/:id", controllers.UpdateBioskop)
	router.DELETE("/bioskop/:id", controllers.DeleteBioskop)

	return router
}