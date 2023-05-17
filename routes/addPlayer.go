package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"veriDrinkApi/game"
)

func AddPlayer(c *gin.Context) {
	sessionId := c.Param("sessionId")
	playerName := c.Param("playerName")

	player := &game.Player{
		Name: playerName,
		Hp:   3,
	}

	sessionManager := game.GetSessionManager()
	session, err := sessionManager.FindSessionById(sessionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error adding player to session",
		})
		return
	}
	if session.Owner != c.ClientIP() {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not the owner of this session",
		})
		return
	}
	session.AddPlayer(player)
	c.JSON(http.StatusOK, gin.H{
		"message": "Player added successfully",
	})
}
