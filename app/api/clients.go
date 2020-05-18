package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nyugoh/sagittarius/app/models"
	"github.com/nyugoh/sagittarius/utils"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	logs := make([]map[string][]models.LogFile, 0)
	for _, folder := range strings.Split(client.Folders, ",") {
		folder = strings.TrimSpace(folder)
		logs = append(logs, map[string][]models.LogFile{
			folder: ListDir(folder, ".log"),
		})
	}

	// Current app logs
	appName := strings.TrimSpace(strings.ToLower(os.Getenv("APP_NAME")))
	logs = append(logs, map[string][]models.LogFile{appName: ListLogs()})

	c.HTML(http.StatusOK, "client.html", gin.H{
		"title":  client.AppName,
		"client": client,
		"logs":   logs,
	})
}

func (app *App) ReadLog(c *gin.Context) {
	logFile := c.Request.URL.Query()["log"]
	utils.Log("Reading log:", logFile)
	if !strings.Contains(logFile[0], ".log") && !strings.Contains(logFile[0], ".sql") {
		_, err := utils.SendMail("admin@quebasetech.co.ke", "Joe Nyugoh", "joenyugoh@gmail.com", "Server Notification", "<p>Someone is trying to access files outside logs folder</p><p>Folder::"+logFile[0]+"</p>")
		if err != nil {
			utils.LogError(err.Error())
		}
		utils.SendError(c, "You are trying to access restricted file: E-mail sent to admin")
		return
	}
	content, err := ioutil.ReadFile(logFile[0])
	if err != nil {
		utils.SendError(c, "Unable to read log file:"+err.Error())
		return
	}
	utils.SendJson(c, gin.H{
		"status":  "success",
		"payload": string(content),
	})
}

func ListLogs() []models.LogFile {
	appName := strings.TrimSpace(strings.ToLower(os.Getenv("APP_NAME")))
	if len(appName) == 0 {
		appName = "app-logs" // Default log file name
	}
	logFolder := os.Getenv(appName + "_LOG_FOLDER")
	logs := ListDir(logFolder, ".log")
	return logs
}

func ListDir(dirPath, fileExt string) []models.LogFile {
	utils.Log("Trying to read ", dirPath, " for ", fileExt, "files")

	files := make([]models.LogFile, 0)
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			utils.LogError(err.Error())
			return err
		}
		if info.IsDir() || filepath.Ext(path) != fileExt {
			return nil
		}
		utils.Log("Path:", path, "Info: Size:", info.Size(), "Name:", info.Name())
		file := models.LogFile{
			Path: path,
			Size: toFixed(float64(info.Size())/1024576.00, 2),
			Date: strings.Split(info.Name(), ".")[1],
			Name: info.Name(),
		}
		files = append(files, file)
		return nil
	})
	if err != nil {
		utils.LogError(err.Error())
		return files
	}
	return files
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
