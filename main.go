package main

import (
	"github.com/gin-gonic/gin"
	"veriDrinkApi/routes"
)

func startResLoading() {
	err := loadFile("res.txt")
	if err != nil {
		return
	}
}

func startRouting() {
	router := gin.Default()

	router.POST("/session/:host", routes.CreateSession)
	router.POST("/session/:sessionId/player/:playerName", routes.AddPlayer)
	router.DELETE("/session/:sessionId/player/:playerName", routes.RemovePlayer)
	router.Run(":8080")
}

func main() {
	startResLoading()
	startRouting()
}
