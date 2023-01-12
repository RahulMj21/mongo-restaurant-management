package main

import (
	"os"

	"github.com/RahulMj21/mongo-restaurant-management/database"
	"github.com/RahulMj21/mongo-restaurant-management/middlewares"
	"github.com/RahulMj21/mongo-restaurant-management/routes"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.client, "food")

func main(){
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app:= gin.New()

	app.Use(gin.Logger())

	routes.UserRoutes(app)
	app.Use(middlewares.Authentication())

	routes.FoodRoutes(app)
	routes.InvoiceRoutes(app)
	routes.MenuRoutes(app)
	routes.OrderItemRoutes(app)
	routes.OrderRoutes(app)
	routes.TableRoutes(app)

	app.Run(":" + port)
}