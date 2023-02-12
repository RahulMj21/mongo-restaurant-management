package controllers

import (
	"context"
	"time"

	"github.com/RahulMj21/mongo-restaurant-management/database"
	"github.com/RahulMj21/mongo-restaurant-management/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func ItemsByOrderId(id string) ([]primitive.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "order_id", Value: id}}}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "food"},
		{Key: "localField", Value: "food_id"},
		{Key: "foreignField", Value: "food_id"},
		{Key: "as", Value: "food"},
	}}}
	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{
		{Key: "path", Value: "$food"},
		{Key: "preserveNullAndEmptyArrays", Value: true},
	}}}

	lookupOrderStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "order"},
		{Key: "localField", Value: "order_id"},
		{Key: "foreignField", Value: "order_id"},
		{Key: "as", Value: "order"},
	}}}
	unwindOrderStage := bson.D{{Key: "$unwind", Value: bson.D{
		{Key: "path", Value: "$order"},
		{Key: "preserveNullAndEmptyArrays", Value: true},
	}}}

	lookupTableStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "table"},
		{Key: "localField", Value: "order.table_id"},
		{Key: "foreignField", Value: "table_id"},
		{Key: "as", Value: "table"},
	}}}
	unwindTableStage := bson.D{{Key: "$unwind", Value: bson.D{
		{Key: "path", Value: "$table"},
		{Key: "preserveNullAndEmptyArrays", Value: true},
	}}}

	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "_id", Value: 0},
		{Key: "amount", Value: "$food.price"},
		{Key: "total_count", Value: 1},
		{Key: "food_name", Value: "$food.name"},
		{Key: "food_image", Value: "$food.image"},
		{Key: "table_number", Value: "$table.table_number"},
		{Key: "price", Value: "$food.price"},
		{Key: "quantity", Value: 1},
	}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: bson.D{
		{Key: "order_id", Value: "$order_id"},
		{Key: "table_id", Value: "$table_id"},
		{Key: "table_number", Value: "$table_number"},
	}},
		{Key: "payment_due", Value: bson.D{{Key: "$sum", Value: "$amount"}}},
		{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
		{Key: "order_items", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
	}}}

	projectStage2 := bson.D{{Key: "$project", Value: bson.D{
		{Key: "_id", Value: 0},
		{Key: "payment_due", Value: 1},
		{Key: "total_count", Value: 1},
		{Key: "table_number", Value: "$_id.table_number"},
		{Key: "order_items", Value: 1},
	}}}

	orderItems := []primitive.M{}
	cursor, err := orderItemCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupStage,
		unwindStage,
		lookupOrderStage,
		unwindOrderStage,
		lookupTableStage,
		unwindTableStage,
		projectStage,
		groupStage,
		projectStage2,
	})

	if err != nil {
		return orderItems, err
	}

	err = cursor.All(ctx, &orderItems)

	return orderItems, err
}

func CreateOrderItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	orderItemPack := OrderItemPack{}
	order := models.Order{}
	if err := c.BindJSON(&orderItemPack); err != nil {
		c.JSON(400, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	order.OrderDate, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.TableId = orderItemPack.TableId
	order_id := OrderItemOrderCreator(order)

	orderItemsToBeInserted := []interface{}{}

	for _, orderItem := range orderItemPack.OrderItems {
		orderItem.OrderId = order_id
		validationErr := validate.Struct(orderItem)
		if validationErr != nil {
			c.JSON(400, gin.H{"status": "fail", "message": validationErr.Error()})
			return
		}
		orderItem.ID = primitive.NewObjectID()
		orderItem.OrderItemId = orderItem.ID.Hex()
		orderItem.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		orderItem.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		num := toFixed(*orderItem.UnitPrice, 2)
		orderItem.UnitPrice = &num

		orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)
	}
	insertedItems, err := orderItemCollection.InsertMany(ctx, orderItemsToBeInserted)
	if err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	c.JSON(201, gin.H{"status": "success", "data": insertedItems})
}

func UpdateOrderItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	orderItemId := c.Param("id")
	orderItem := models.OrderItem{}

	if orderItemId == "" {
		c.JSON(400, gin.H{"status": "fail", "message": "id cannot be empty"})
		return
	}
	if err := c.BindJSON(&orderItem); err != nil {
		c.JSON(400, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var orderItemObj primitive.D

	if orderItem.UnitPrice != nil {
		unitPrice := toFixed(*orderItem.UnitPrice, 2)
		orderItemObj = append(orderItemObj, bson.E{Key: "unit_price", Value: &unitPrice})
	}
	if orderItem.FoodId != nil {
		orderItemObj = append(orderItemObj, bson.E{Key: "food_id", Value: orderItem.FoodId})
	}
	if orderItem.Quantity != nil {
		orderItemObj = append(orderItemObj, bson.E{Key: "quantity", Value: orderItem.Quantity})
	}

	orderItem.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	orderItemObj = append(orderItemObj, bson.E{Key: "updated_at", Value: orderItem.UpdatedAt})

	filter := bson.D{{Key: "order_item_id", Value: orderItemId}}
	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	result, err := orderItemCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: orderItemObj}}, &opt)
	if err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "fail", "data": result})
}
