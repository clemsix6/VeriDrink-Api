package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"veriDrinkApi/game"
)

func RemovePlayer(c *gin.Context) {
	sessionId := c.Param("sessionId")
	playerName := c.Param("playerName")

	sessionManager := game.GetSessionManager()
	session, err := sessionManager.FindSessionById(sessionId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Session not found",
		})
		return
	}
	if session.Owner != c.ClientIP() {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not the owner of this session",
		})
		return
	}
	err = session.RemovePlayer(playerName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error removing player from session",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Player removed successfully",
	})
}
