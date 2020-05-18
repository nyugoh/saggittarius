package api

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	io "github.com/googollee/go-socket.io"
	"github.com/jinzhu/gorm"
	"net/http"
	"os"
	"time"
)

type App struct {
	DB    *gorm.DB
	Name  string
	Port  string
	Redis redis.Client
	SocketServer *io.Server
}

func (app *App) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{"title": "Dashboard"})
}

func (app *App) GenerateToken(appName string, hash string) (token string, err error) {
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["app_name"] = appName
	atClaims["app_hash"] = hash
	atClaims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err = at.SignedString([]byte(os.Getenv("SESSION_KEY")))
	if err != nil {
		return "", err
	}
	return token, nil
}
