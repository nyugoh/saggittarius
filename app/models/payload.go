package models

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	AppName string `json:"app_name"`
	jwt.StandardClaims
}

type LogFile struct {
	Name string  `json:"name"`
	Size float64 `json:"size"`
	Path string  `json:"path"`
	Date string  `json:"date"`
}