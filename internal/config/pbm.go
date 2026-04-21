package config

import (
	"fmt"
	"io"
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

type PbPath struct {
	Local string `yaml:"local"`

	Remote string `yaml:"remote"`
	Ref    string `yaml:"ref"`
	File   string `yaml:"file"`
}

type Input struct {
	PbPath  `yaml:",inline"`
	DescOut string `yaml:"desc_out"`
}

type Dep struct {
	PbPath `yaml:",inline"`
}

func PbmConfigFromFile(path string) (*PbmConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return PbmConfigFromReader(file)
}

func PbmConfigFromReader(reader io.Reader) (*PbmConfig, error) {
	configBytes, err := io.ReadAll(reader)
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
	if len(config.Input) == 0 || len(config.Gen) == 0 {
		return nil, fmt.Errorf("input or gen is invalid")
	}
	return &config, nil
}

func InitConfig() error {
	_, err := os.Stat("pbm.yaml")
	if err == nil {
		return fmt.Errorf("pbm.yaml already exist")
	}
	examplePbmConfig := `version: v1
deps:
  - remote: https://github.com/pbm-org/pbm.git
    ref: main
  - remote: git@github.com:pbm-org/pbm.git
    ref: v0.0.1
  - local: proto1/pbm.proto
gen:
  - plugin: go
    out: .
    opt:
      - paths=source_relative
input:
  - local: proto/proto1.proto
    desc_out: ./xxx/xxx
`
	err = os.WriteFile("pbm.yaml", []byte(examplePbmConfig), 0755)
	return err
}
