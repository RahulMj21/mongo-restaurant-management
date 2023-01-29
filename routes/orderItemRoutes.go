package routes

import (
	"github.com/RahulMj21/mongo-restaurant-management/controllers"
	"github.com/gin-gonic/gin"
)

func OrderItemRoutes(api *gin.RouterGroup) {
	api.GET("/order-items", controllers.GetOrderItems)
	api.GET("/order-items/:id", controllers.GetOrderItem)
	api.GET("/order-items-order/:order_id", controllers.GetOrderItemsByOrderId)
	api.POST("/order-items", controllers.CreateOrderItem)
	api.PATCH("/order-items/:id", controllers.UpdateOrderItem)
}
