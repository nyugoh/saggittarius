package api

import (
	"fmt"
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
		return
	}

	utils.SendJson(c, gin.H{
		"success": "success",
		"payload": clients,
	})
}

func (app *App) GetFolders(c *gin.Context) {
	appName := utils.ExtractAppName(c)
	if len(appName) == 0 {
		utils.SendError(c, "Unable to read client name")
		return
	}
	client := models.Client{}
	if err := app.DB.Find(&client, "app_name=?", appName).Error; err != nil {
		utils.SendError(c, err.Error())
		return
	}

	utils.SendJson(c, gin.H{
		"success": "success",
		"folders": client.Folders,
		"config":  client.ConfigPath,
	})
}

func (app *App) ListFolders(c *gin.Context) {
	clientId := c.Param("id")
	client := models.Client{}
	if err := app.DB.Find(&client, "id=?", clientId).Error; err != nil {
		utils.SendError(c, err.Error())
		return
	}
	url := fmt.Sprintf("http://%s:%s/logs", client.IPAddress, client.Port)
	res, err := utils.SendGet(url)
	if err != nil {
		utils.SendError(c, err.Error())
		return
	}

	c.HTML(http.StatusOK, "client.html", gin.H{
		"title":  client.AppName,
		"client": client,
		"logs":   res["logs"],
	})
}

func (app *App) DeleteClient(c *gin.Context) {
	clientId := c.Param("id")
	client := models.Client{}
	if err := app.DB.Find(&client, "id=?", clientId).Error; err != nil {
		utils.SendError(c, err.Error())
		return
	}

	utils.Log("Deleting a client:", client.AppName)
	app.DB.Delete(&client)

	utils.SendJson(c, gin.H{
		"success": "success",
		"message": "client deleted successfully",
	})
}

func (app *App) EditClient(c *gin.Context) {
	payload := struct {
		AppId    string `json:"appId" `
		AppName    string `json:"appName" `
		IPAddress  string `json:"appIp"`
		Port       string `json:"appPort"`
		Folders    string `json:"folders"`
		ConfigPath string `json:"configPath"`
	}{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.SendError(c, err.Error())
		return
	}

	clientId := c.Param("id")
	client := models.Client{}
	if err := app.DB.Find(&client, "id=?", clientId).Error; err != nil {
		utils.SendError(c, err.Error())
		return
	}

	utils.Log("Updating a client:", payload.AppName)
	client.AppName = payload.AppName
	client.IPAddress = payload.IPAddress
	client.Port = payload.Port
	client.Folders = payload.Folders
	client.ConfigPath = payload.ConfigPath
	app.DB.Save(&client)

	utils.SendJson(c, gin.H{
		"success": "success",
		"message": "client updated successfully",
		"client": client,
	})
}
