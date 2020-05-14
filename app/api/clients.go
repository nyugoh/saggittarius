package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nyugoh/sagittarius/app/models"
	"github.com/nyugoh/sagittarius/utils"
	"net/http"
)

func (app *App) Clients(c *gin.Context) {
	c.HTML(http.StatusOK, "clients.html", gin.H{
		"title": "Clients",
	})
}

func (app *App) AddClient(c *gin.Context) {
	payload := struct {
		AppName    string `json:"appName"`
		AppIP      string `json:"appIp"`
		AppPort    string `json:"appPort"`
		Folders    string `json:"folders"`
		ConfigPath string `json:"configPath"`
	}{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		utils.SendError(c, "Error:"+err.Error())
		return
	}

	utils.Log("Received a request to create a new client:", payload)

	client := models.Client{
		AppName:    payload.AppName,
		IPAddress:  payload.AppIP,
		Port:       payload.AppPort,
		Folders:    payload.Folders,
		ConfigPath: payload.ConfigPath,
		ClientHash: utils.GenerateUUID(),
		Status:     true,
	}

	if err := app.DB.Save(&client).Error; err != nil {
		utils.SendError(c, "Error creating client:"+err.Error())
		return
	}

	utils.Log("Client created successfully:", client)
	utils.SendJson(c, gin.H{
		"success": "success",
		"payload": client,
	})
}

func (app *App) ListClients(c *gin.Context) {
    utils.Log("Fetching list of clients")

    clients := make([]models.Client, 0)
    if err := app.DB.Find(&clients).Error; err != nil {
    	utils.SendError(c, err.Error())
	}

	utils.SendJson(c, gin.H{
		"success": "success",
		"payload": clients,
	})
}
