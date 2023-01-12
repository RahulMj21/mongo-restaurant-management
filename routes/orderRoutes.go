package routes

import (
	"github.com/RahulMj21/mongo-restaurant-management/controllers"
	"github.com/gin-gonic/gin"
)

func OrderRoutes(app *gin.Engine) {
	app.GET("/orders",controllers.GetOrders)
	app.GET("/order/:id",controllers.GetOrder)
	app.POST("/orders",controllers.CreateOrder)
	app.PATCH("/orders/:id",controllers.UpdateOrder)
}