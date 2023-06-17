package config

import (
	"dotkafx/log"
	"dotkafx/model"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const configFile = "dotkafx_config.yml"

// GetConfigData will check the Home folder, of the user running the application, for the dotkafx_config.yml file,
// and return its content. If the file cannot be found (first time run) then it will be created with the defaultConfig data.
func GetConfigData(defaultConfig []byte) ([]byte, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Error("Failed to find Home folder: %s", err)
		return defaultConfig, nil
	}

	configFilePath := filepath.Join(homeDir, configFile)
	if _, err := os.Stat(configFilePath); err == nil {
		return os.ReadFile(configFilePath)
	} else {
		log.Warn("Failed to read config file in Home folder: %s", err)
	}

	if err := os.WriteFile(configFilePath, defaultConfig, 0644); err != nil {
		return nil, err
	}

	return defaultConfig, nil
}

// CreateConfig accepts the content of a file as argument, and creates the application configuration object from it.
func CreateConfig(configData []byte) (model.Config, error) {
	var inputConf model.ConfigInput

	if err := yaml.Unmarshal(configData, &inputConf); err != nil {
		return model.Config{}, err
	}

	return inputConf.Parse()
}
