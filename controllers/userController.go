package controllers

import "github.com/gin-gonic/gin"

func GetUsers(c * gin.Context){
	c.JSON(200, "hello")
}

func GetUser(c *gin.Context) {
	c.JSON(200, "hello")
}

func SignUp(c *gin.Context) {
	c.JSON(200, "hello")
}

func Login(c *gin.Context) {
	c.JSON(200, "hello")
}

func HashPassword(password string) string{
	return password
}

func VerifyPassword(hashedPassword string, password string) (bool, string){
	return password==hashedPassword,"hello"
}