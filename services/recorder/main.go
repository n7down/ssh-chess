package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func playerAuthHandler(c *gin.Context) {
	a := AuthPlayerRequest{}
	err := c.BindJSON(&a)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error with binding json: %v", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	secretIsValid, err := a.CheckAuth()
	if err != nil {
		fmt.Println(fmt.Sprintf("Error with secret is valid: %v", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": secretIsValid})
}

func createNewPlayerHandler(c *gin.Context) {
	p := CreatePlayerRequest{}
	err := c.BindJSON(&p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	_, err = p.CreateNewPlayer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success!"})
}

func checkUserNameHandler(c *gin.Context) {
	p := CheckUserNameRequest{}
	err := c.BindJSON(&p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	userNameExists, err := p.CheckUserNameExists()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": userNameExists})
}

func main() {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("/player/create", createNewPlayerHandler)
		api.GET("/player/checkuser", checkUserNameHandler)
		api.GET("/player/auth", playerAuthHandler)
	}

	ws := router.Group("/ws")
	{
		ws.GET("/record", func(c *gin.Context) {
			wsHandler(c.Writer, c.Request)
		})
	}

	router.Run(":8000")
}
