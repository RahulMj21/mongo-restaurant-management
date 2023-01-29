package routes

import (
	"github.com/RahulMj21/mongo-restaurant-management/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(api *gin.RouterGroup) {
	api.GET("/users", controllers.GetUsers)
	api.GET("/users/:id", controllers.GetUser)
	api.GET("/signup", controllers.SignUp)
	api.GET("/login", controllers.Login)
}
