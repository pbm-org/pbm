package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type PbmConfig struct {
	Version string  `yaml:"version"`
	Gen     []Gen   `yaml:"gen"`
	Input   []Input `yaml:"input"`
	Deps    []Dep   `yaml:"deps"`
}

type Gen struct {
	Plugin string   `yaml:"plugin"`
	Out    string   `yaml:"out"`
	Opt    []string `yaml:"opt"`
}

type Input struct {
	Path string `yaml:"path"`
	Dir  string `yaml:"dir"`
	File string `yaml:"file"`
}

type Dep struct {
	Path string `yaml:"path"`
	Ref  string `yaml:"ref"`
}

func PbmConfigFromFile(path string) (*PbmConfig, error) {
	configBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config = PbmConfig{
		Deps:  []Dep{},
		Gen:   []Gen{},
		Input: []Input{},
	}
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
