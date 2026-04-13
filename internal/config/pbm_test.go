package config

import (
	"reflect"
	"testing"
)

func TestConfig(t *testing.T) {
	pbmConfig, err := PbmConfigFromFile("../../pbm.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if pbmConfig.Version != "v1" {
		t.Errorf("expected version %s, got %s", "", pbmConfig.Version)
	}
	deps := []Dep{{
		Path: "https://github.com/pbm-org/pbm.git",
		Ref:  "main",
	},
		{
			Path: "git@github.com:pbm-org/pbm.git",
			Ref:  "v0.0.1",
		},
		{
			Path: "proto1/pbm.proto",
		},
		{
			Path: "proto2/pbm.proto",
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
			Path: "proto/proto1.proto",
		},
		{
			Path: "git@github.com:labulakalia/pbb.git",
			Dir:  "proto",
			File: "proto/proto2.proto",
		},
		{
			Path: "proto_dir",
		},
	}
	if !reflect.DeepEqual(pbmConfig.Input, input) {
		t.Errorf("expected input %v, got %v", input, pbmConfig.Input)
	}
}
