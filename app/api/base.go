package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	io "github.com/googollee/go-socket.io"
	"github.com/jinzhu/gorm"
	"net/http"
)

type App struct {
	DB    *gorm.DB
	Name  string
	Port  string
	Redis redis.Client
	SocketServer *io.Server
}

func (app *App) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{"title": "Sagittarius"})
}