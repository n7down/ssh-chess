package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllActivePlayersHandler(c *gin.Context) {

}

func GetAllPlayersHandler(c *gin.Context) {

	type AllActivePlayers struct {
		players []Player
	}
}

func CreateNewPlayerHandler(c *gin.Context) {
	p := Player{}
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

func CheckSecretHandler(c *gin.Context) {
	req := CheckSecretRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	secretIsValid, err := req.CheckHashedSecret()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": secretIsValid})
}

func GetAllActiveGamesHandler(c *gin.Context) {

}

func GetAllGamesHandler(c *gin.Context) {

}
