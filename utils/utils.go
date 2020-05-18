package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	uuid2 "github.com/google/uuid"
	"github.com/nyugoh/sagittarius/app/models"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

func SendError(c *gin.Context, msg string)  {
	LogError(msg)
	c.JSON(http.StatusBadRequest, gin.H{
		"error": msg,
	})
}

func SendJson(c *gin.Context, payload gin.H)  {
	c.JSON(http.StatusOK, payload)
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	RespondWithJSON(w, code, map[string]interface{}{"error": msg})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	Log(fmt.Sprintf("RESPONSE:: Status:%d Payload: %v", code, payload))
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func CurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func ValidateEmail(email string) (bool, error) {
	var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(email) > 254 || !rxEmail.MatchString(email) {
		return false, errors.New("email is invalid")
	}
	return true, nil
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SESSION_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}


func AuthRequired() gin.HandlerFunc  {
	return func(c *gin.Context) {
		Log("Hit an auth required endpoint...")
		auth := c.Request.Header.Get("Authorization")
		if len(auth)== 0 {
			SendError(c, "Authorization token is required.")
			return
		}
		token, err := VerifyToken(auth)
		if err != nil {
			SendError(c, err.Error())
			return
		}

		claims := token.Claims.(*models.CustomClaims)
		Log(claims.AppName)
		Log("Request made by app name:", claims.AppName)
		c.Set("app_name", claims.AppName)
		c.Next()
	}
}

func MetricsMonitor() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		// before request
		c.Next()
		// after request
		latency := time.Since(t)
		// access the status we are sending
		status := c.Writer.Status()
		if !strings.Contains(c.Request.URL.Path, "socket.io"){ // Avoid logging for /socket.io
			log.Println("Request: ", c.Request.URL.Path,  "took:", latency, "Status:", status)
		}
	}
}

func ExtractAppName(c *gin.Context) (appName string) {
	if name, ok := c.Get("app_name"); ok {
		appName = fmt.Sprintf("%v", name)
	} else {
		appName = "Unknown"
	}

	return
}

func SendGet(url string) (result map[string]interface{}, err error) {
	client := http.Client{}
	request, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer([]byte("")))

	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	json.NewDecoder(res.Body).Decode(&result)

	Log(result)
	if res.StatusCode == 200 {
		return result, nil
	} else {
		errMsg := fmt.Sprintf("%s", result["error"])
		LogError(errMsg)
		return nil, errors.New(errMsg)
	}
}

func ExitApp(code int) {
	Log("Exiting app...")
	os.Exit(code)
}

func GenerateUUID() string {
	uuid := uuid2.New()
	return uuid.String()
}