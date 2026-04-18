package config

import (
	"reflect"
	"strings"
	"testing"
)

func TestConfig(t *testing.T) {
	pbmConfigStr := `version: v1
deps:
  - remote: https://github.com/pbm-org/pbm.git
    ref: main
  - remote: git@github.com:pbm-org/pbm.git
    ref: v0.0.1
  - local: proto1/pbm.proto
  - local: proto2/pbm.proto
gen:
  - plugin: go
    out: .
    opt:
      - paths=source_relative
  - plugin: dart
    out: ./gen_dart
input:
  - local: proto/proto1.proto
    desc_out: ./xxx/xxx
  - remote: git@github.com:labulakalia/pbb.git
    file: proto/proto.proto
  - local: proto_dir
`
	pbmConfig, err := PbmConfigFromReader(strings.NewReader(pbmConfigStr))
	if err != nil {
		t.Fatal(err)
	}
	if pbmConfig.Version != "v1" {
		t.Errorf("expected version %s, got %s", "", pbmConfig.Version)
	}
	deps := []Dep{{
		PbPath: PbPath{
			Remote: "https://github.com/pbm-org/pbm.git",
			Ref:    "main",
		},
	},
		{
			PbPath: PbPath{
				Remote: "git@github.com:pbm-org/pbm.git",
				Ref:    "v0.0.1",
			},
		},
		{
			PbPath: PbPath{
				Local: "proto1/pbm.proto",
			},
		},
		{
			PbPath: PbPath{
				Local: "proto2/pbm.proto",
			},
		}}
	if !reflect.DeepEqual(pbmConfig.Deps, deps) {
		t.Errorf("expected deps %v, got %v", deps, pbmConfig.Deps)
	}

	gens := []Gen{
		{
			Plugin: "go",
			Out:    ".",
			Opt:    []string{"paths=source_relative"},
		},
		{
			Plugin: "dart",
			Out:    "./gen_dart",
		},
	}
	if !reflect.DeepEqual(pbmConfig.Gen, gens) {
		t.Errorf("expected gens %v, got %v", gens, pbmConfig.Gen)
	}
	input := []Input{
		{
			PbPath: PbPath{
				Local: "proto/proto1.proto",
			},
			DescOut: "./xxx/xxx",
		},
		{
			PbPath: PbPath{Remote: "git@github.com:labulakalia/pbb.git", File: "proto/proto.proto"},
		},
		{
			PbPath: PbPath{
				Local: "proto_dir",
			},
		},
	}
	if !reflect.DeepEqual(pbmConfig.Input, input) {
		t.Errorf("expected input \n%+v, got \n%+v", input, pbmConfig.Input)
	}
}
