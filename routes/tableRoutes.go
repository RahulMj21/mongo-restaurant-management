package routes

import (
	"github.com/RahulMj21/mongo-restaurant-management/controllers"
	"github.com/gin-gonic/gin"
)

func TableRoutes(api *gin.RouterGroup) {
	api.GET("/tables", controllers.GetTables)
	api.GET("/tables/:id", controllers.GetTable)
	api.POST("/tables", controllers.CreateTable)
	api.PATCH("/tables/:id", controllers.UpdateTable)
}
