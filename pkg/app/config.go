package app

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Replacement struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}

type Endpoint struct {
	Path         string        `yaml:"path"`
	TargetUrl    string        `yaml:"targetUrl"`
	Replacements []Replacement `yaml:"replacements"`
}

func (e *Endpoint) getTimeout() time.Duration {
	d, _ := time.ParseDuration("5m0s")

	return d
}

type Config struct {
	Endpoints []Endpoint `yaml:"endpoints"`
}

// LoadConfigurationFromFile is reading a YAML configuration file into internal structs
func LoadConfigurationFromFile(path string) (Config, error) {
	f, readErr := os.ReadFile(path)

	if readErr != nil {
		return Config{}, fmt.Errorf("cannot read file at path '%s', error: '%s'", path, readErr.Error())
	}

	cfg := Config{}
	parseErr := yaml.Unmarshal(f, &cfg)
	if parseErr != nil {
		return Config{}, fmt.Errorf("cannot parse yaml at path '%s', error: '%s'", path, parseErr.Error())
	}

	return cfg, nil
}
