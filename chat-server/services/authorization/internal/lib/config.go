package lib

import (
	"fmt"
	"os"

	"github.com/go-yaml/yaml"
	"github.com/gox7/notify/services/authorization/models"
)

func NewConfig(model *models.LocalConfig) {
	configPath := "config/config.yaml"

	// Read the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println("[-] config.read:", err)
		os.Exit(1)
	}

	// Parse YAML into the model
	if err := yaml.Unmarshal(data, model); err != nil {
		fmt.Println("[-] config.unmarshal:", err)
		os.Exit(1)
	}
}
