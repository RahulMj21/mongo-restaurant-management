package routes

import (
	"github.com/RahulMj21/mongo-restaurant-management/controllers"
	"github.com/gin-gonic/gin"
)

func MenuRoutes(api *gin.RouterGroup) {
	api.GET("/menus", controllers.GetMenus)
	api.GET("/menus/:id", controllers.GetMenu)
	api.POST("/menus", controllers.CreateMenu)
	api.PATCH("/menus/:id", controllers.UpdateMenu)
}
