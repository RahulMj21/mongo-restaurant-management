package controllers

import (
	"context"
	"time"

	"github.com/RahulMj21/mongo-restaurant-management/database"
	"github.com/RahulMj21/mongo-restaurant-management/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItemPack struct {
	TableId    *string
	OrderItems []models.OrderItem
}

var orderItemCollection = database.OpenCollection(database.Client, "order_item")

func GetOrderItems(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	cursor, err := orderItemCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	allOrderItems := []primitive.M{}
	if err := cursor.All(ctx, &allOrderItems); err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "success", "data": allOrderItems})
}

func GetOrderItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	orderItemId := c.Param("id")
	if orderItemId == "" {
		c.JSON(400, gin.H{"status": "fail", "message": "order_item id cannot be empty"})
		return
	}

	orderItem := models.OrderItem{}

	err := orderItemCollection.FindOne(ctx, bson.D{{Key: "order_item_id", Value: orderItemId}}).Decode(&orderItem)
	if err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": "cannot find the order_item"})
		return
	}

	c.JSON(200, gin.H{"status": "success", "data": orderItem})
}

func GetOrderItemsByOrderId(c *gin.Context) {
	orderID := c.Param("order_id")
	if orderID == "" {
		c.JSON(400, gin.H{"status": "fail", "message": "order id cannot be empty"})
		return
	}
	allOrderItems, err := ItemsByOrderId(orderID)
	if err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": "cannot get order items"})
		return
	}

	c.JSON(200, gin.H{"status": "success", "data": allOrderItems})
}

func ItemsByOrderId(id string) (OrderItems []primitive.M, err error) {
	return
}

func CreateOrderItem(c *gin.Context) {
	_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	orderItemPack := OrderItemPack{}
	// order := models.Order{}
	if err := c.BindJSON(&orderItemPack); err != nil {
		c.JSON(400, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	c.JSON(201, gin.H{"status": "success", "message": "order created"})
}

func UpdateOrderItem(c *gin.Context) {
	_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	c.JSON(200, gin.H{"status": "fail", "message": "orderItemUpdated"})
}
