package services

import (
	"gopkg.in/ini.v1"
	"os"
)

func GetProperty(sectionName string, keyName string) string {
	data, _ := os.ReadFile("../config.ini")
	cfg, _ := ini.Load(data)

	section, _ := cfg.GetSection(sectionName)
	k, _ := section.GetKey(keyName)
	return k.String()
}
