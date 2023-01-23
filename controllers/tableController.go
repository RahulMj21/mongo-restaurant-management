package controllers

import (
	"github.com/RahulMj21/mongo-restaurant-management/database"
	"github.com/gin-gonic/gin"
)

var tableCollection = database.OpenCollection(database.Client, "table")

func GetTables(c *gin.Context) {

}

func GetTable(c *gin.Context) {

}

func CreateTable(c *gin.Context) {

}

func UpdateTable(c *gin.Context) {

}
