package routes

import (
	"github.com/RahulMj21/mongo-restaurant-management/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(app *gin.Engine) {
	app.GET("/users", controllers.GetUsers)
	app.GET("/users/:id", controllers.GetUser)
	app.GET("/signup", controllers.SignUp)
	app.GET("/login", controllers.Login)
}
