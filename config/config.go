package config

import (
	"dotkafx/log"
	"dotkafx/model"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const configFile = "dotkafx_config.yml"

func GetConfigData(defaultConfig []byte) ([]byte, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Error("Failed to find Home folder: %s", err)
		return defaultConfig, nil
	}

	configFilePath := filepath.Join(homeDir, configFile)
	if _, err := os.Stat(configFilePath); err == nil {
		return ioutil.ReadFile(configFilePath)
	} else {
		log.Warn("Failed to read config file in Home folder: %s", err)
	}

	if err := ioutil.WriteFile(configFilePath, []byte(defaultConfig), 0644); err != nil {
		return nil, err
	}

	return defaultConfig, nil
}

func CreateConfig(configData []byte) (model.Config, error) {
	var inputConf model.ConfigInput

	if err := yaml.Unmarshal(configData, &inputConf); err != nil {
		return model.Config{}, err
	}

	return inputConf.Parse()
}
