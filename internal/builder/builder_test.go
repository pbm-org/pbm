package builder

import (
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/pbm-org/pbm/internal/config"
)

func TestBuildProto(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	pbmConfigStr := `version: v1
deps:
  - remote: https://cnb.cool/medianexapp/plugin_api
    ref: main
  - remote: https://github.com/googleapis/googleapis
    ref: master
  - remote: https://github.com/bufbuild/protoc-gen-validate
    ref: main
gen:
  - plugin: go-lite
    out: .
    opt:
      - paths=source_relative
  - plugin: validate-go
    out: .
    opt:
      - paths=source_relative
input:
  - local: testdata/proto
    desc_out: testdata/proto/proto1.pb
`
	os.Chdir("../../")
	pbmConfig, err := config.PbmConfigFromReader(strings.NewReader(pbmConfigStr))
	if err != nil {
		t.Fatal(err)
	}
	err = CheckPbCfg(pbmConfig)
	if err != nil {
		t.Fatal(err)
	}

	cmds, err := PbBuildCmd(pbmConfig)
	if err != nil {
		t.Fatal(err)
	}
	if len(cmds) == 0 {
		t.Fatal("cmds is 0")
	}
	err = RunPbmCmd(cmds)
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Stat("testdata/proto/proto1.pb")
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Stat("testdata/proto/proto1.pb.go")
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Stat("testdata/proto/proto1.pb.go")
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Stat("testdata/proto/proto1.pb.validate.go")
	if err != nil {
		t.Fatal(err)
	}
	os.Remove("testdata/proto/proto1.pb")
	os.Remove("testdata/proto/proto1.pb.go")
	os.Remove("testdata/proto/proto1.pb.validate.go")
}
