package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nyugoh/sagittarius-client/utils"
	"github.com/nyugoh/sagittarius/app/models"
)

func (app *App) GenerateJWT(c *gin.Context) {
    payload := struct {
		AppName string `json:"appName"`
		AppHash string `json:"appHash"`
		AppIp string `json:"appIp"`
		AppPort string `json:"appPort"`
	}{}
    if err := c.ShouldBindJSON(&payload); err != nil {
    	utils.SendError(c, err.Error())
		return
	}

	utils.Log("Received a request to generate token:", payload)

    client := models.Client{}

    // Fetch client
    if err := app.DB.Find(&client, "app_name=?", payload.AppName).Error; err != nil {
    	utils.SendError(c, err.Error())
		return
	}

	token, err := app.GenerateToken(payload.AppName, payload.AppHash)
	if err != nil {
		utils.SendError(c, err.Error())
		return
	}

	utils.SendJson(c, gin.H{
		"success": "success",
		"payload": token,
	})
}
