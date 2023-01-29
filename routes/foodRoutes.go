package routes

import (
	"github.com/RahulMj21/mongo-restaurant-management/controllers"
	"github.com/gin-gonic/gin"
)

func FoodRoutes(api *gin.RouterGroup) {
	api.GET("/foods", controllers.GetFoods)
	api.GET("/foods/:id", controllers.GetFood)
	api.POST("/foods", controllers.CreateFood)
	api.PATCH("/foods/:id", controllers.UpdateFood)
}
