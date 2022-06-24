package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

func Load(configPath string) (*Config, error) {
	config := &Config{}
	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}
