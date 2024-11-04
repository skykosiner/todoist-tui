package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	Token string `json:"token"`
}

func NewConfig() (Config, error) {
	var config Config
	bytes, err := os.ReadFile(fmt.Sprintf("%s/.config/todoist/settings.json", os.Getenv("HOME")))
	if err != nil {
		return config, errors.New("Config doesn't eixist. Check the README.")
	}

	if err := json.Unmarshal(bytes, &config); err != nil {
		return config, err
	}

	return config, nil
}
