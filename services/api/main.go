package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.GET("/games/active", GetAllActiveGamesHandler)
		api.GET("/games", GetAllGamesHandler)
		api.GET("/players/active", GetAllActivePlayersHandler)
		api.GET("/players", GetAllPlayersHandler)
		api.POST("/player/create", CreateNewPlayerHandler)

		api.GET("/player/secret", CheckSecretHandler)
	}

	router.Run(":8080")
}
