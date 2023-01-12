package routes

import (
	"github.com/RahulMj21/mongo-restaurant-management/controllers"
	"github.com/gin-gonic/gin"
)

func TableRoutes(app *gin.Engine) {
	app.GET("/tables",controllers.GetTables)
	app.GET("/table/:id",controllers.GetTable)
	app.POST("/tables",controllers.CreateTable)
	app.PATCH("/tables/:id",controllers.UpdateTable)
}