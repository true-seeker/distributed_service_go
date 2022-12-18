package services

import (
	"gopkg.in/ini.v1"
	"os"
)

// GetProperty Получение проперти из конфига
func GetProperty(sectionName string, keyName string) string {
	data, err := os.ReadFile("../config.ini")
	if err != nil {
		data, err = os.ReadFile("config.ini")
		if err != nil {
			FailOnError(err, "cant open config")
		}
	}
	cfg, _ := ini.Load(data)

	section, _ := cfg.GetSection(sectionName)
	k, _ := section.GetKey(keyName)
	return k.String()
}
