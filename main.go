package main

import (
	"os"

	"github.com/RahulMj21/mongo-restaurant-management/middlewares"
	"github.com/RahulMj21/mongo-restaurant-management/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app := gin.New()
	api := app.Group("/api/v1")

	api.Use(gin.Logger())

	routes.UserRoutes(api)
	api.Use(middlewares.Authentication)

	routes.FoodRoutes(api)
	routes.InvoiceRoutes(api)
	routes.MenuRoutes(api)
	routes.OrderItemRoutes(api)
	routes.OrderRoutes(api)
	routes.TableRoutes(api)

	app.Run(":" + port)
}
