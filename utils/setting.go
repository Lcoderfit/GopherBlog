package utils

import (
	"fmt"
	"gopkg.in/ini.v1"
)

var (
	AppMode string
	HttpPort string

	DatabaseType string
	DatabaseName string
	Host          string
	Port          string
	Username      string
	Password      string
)


func init() {
	file, err := ini.Load("../config/config.ini")
	if err != nil {
		fmt.Println()
	}
}