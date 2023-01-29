package routes

import (
	"github.com/RahulMj21/mongo-restaurant-management/controllers"
	"github.com/gin-gonic/gin"
)

func OrderRoutes(api *gin.RouterGroup) {
	api.GET("/orders", controllers.GetOrders)
	api.GET("/orders/:id", controllers.GetOrder)
	api.POST("/orders", controllers.CreateOrder)
	api.PATCH("/orders/:id", controllers.UpdateOrder)
}
