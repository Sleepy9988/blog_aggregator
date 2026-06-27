package config

import (
	"encoding/json"
	"os"
)

func ReadJson() (Config, error) {
	file_path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	jsonFile, err := os.Open(file_path)
	if err != nil {
		return Config{}, err
	}

	defer jsonFile.Close()

	decoder := json.NewDecoder(jsonFile)

	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
