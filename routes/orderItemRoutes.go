package routes

import (
	"github.com/RahulMj21/mongo-restaurant-management/controllers"
	"github.com/gin-gonic/gin"
)

func OrderItemRoutes(app *gin.Engine) {
	app.GET("/order-items", controllers.GetOrderItems)
	app.GET("/order-items/:id", controllers.GetOrderItem)
	app.GET("/order-items-order/:order_id", controllers.GetOrderItemsByOrderId)
	app.POST("/order-items", controllers.CreateOrderItem)
	app.PATCH("/order-items/:id", controllers.UpdateOrderItem)
}
