package routes

import (
	"github.com/RahulMj21/mongo-restaurant-management/controllers"
	"github.com/gin-gonic/gin"
)

func FoodRoutes(app *gin.Engine) {
	app.GET("/foods", controllers.GetFoods)
	app.GET("/foods/:id", controllers.GetFood)
	app.POST("/foods", controllers.CreateFood)
	app.PATCH("/foods/:id", controllers.UpdateFood)
}
