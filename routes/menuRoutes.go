package routes

import (
	"github.com/RahulMj21/mongo-restaurant-management/controllers"
	"github.com/gin-gonic/gin"
)

func MenuRoutes(app *gin.Engine) {
	app.GET("/menus",controllers.GetMenus)
	app.GET("/menu/:id",controllers.GetMenu)
	app.POST("/menus",controllers.CreateMenu)
	app.PATCH("/menus/:id",controllers.UpdateMenu)
}