package routes

import (
	"github.com/RahulMj21/mongo-restaurant-management/controllers"
	"github.com/gin-gonic/gin"
)

func InvoiceRoutes(api *gin.RouterGroup) {
	api.GET("/invoices", controllers.GetInvoices)
	api.GET("/invoices/:id", controllers.GetInvoice)
	api.POST("/invoices", controllers.CreateInvoice)
	api.PATCH("/invoices/:id", controllers.UpdateInvoice)

}
