package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"veriDrinkApi/game"
)

func CreateSession(c *gin.Context) {

	owner := c.ClientIP()
	session, err := game.GetSessionManager().NewSession(owner)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la création de la session",
		})
		return
	}

	c.JSON(http.StatusOK, session)
}
