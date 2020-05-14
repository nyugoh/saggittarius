package models

import "github.com/jinzhu/gorm"

type Client struct {
	gorm.Model
	AppName    string `json:"appName"  gorm:"not null; unique"`
	IPAddress  string `json:"ipAddress" gorm:"not null;"`
	Port       string `json:"port" gorm:"not null;"`
	Folders    string `json:"folders"`
	ConfigPath string `json:"configPath"`
	ClientHash string `json:"clientHash"`
	AuthToken  string `json:"authToken"`
	Status     bool   `json:"status" gorm:"default: 1"`
}
