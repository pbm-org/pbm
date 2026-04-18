package builder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pbm-org/pbm/internal/config"
)

func TestBuilder(t *testing.T) {
	pbmConfigStr := `version: v1
deps:
  - remote: https://cnb.cool/medianexapp/plugin_api
    ref: main
  - remote: git@github.com:pbm-org/pbm.git
    ref: v0.0.1
  - local: proto1
  - local: proto2
gen:
  - plugin: go
    out: .
    opt:
      - paths=source_relative
  - plugin: dart
    out: ./gen_dart
input:
  - local: proto/proto1.proto
    desc_out: ./xxx/xxx.proto
  - remote: git@github.com:labulakalia/pbb.git
    ref: main
    file: proto/*.proto
  - local: proto_dir/*.proto
`
	pbmConfig, err := config.PbmConfigFromReader(strings.NewReader(pbmConfigStr))
	if err != nil {
		t.Fatal(err)
	}
	cmds, err := PbBuildCmd(pbmConfig)
	if err != nil {
		t.Fatal(err)
	}
	for _, cmd := range cmds {
		fmt.Println(cmd)
	}
}

func TestGolb(t *testing.T) {
	fmt.Println(filepath.Glob("wde*dwe/builder.go"))
}

func TestBuildProto(t *testing.T) {
	pbmConfigStr := `version: v1
deps:
  - remote: https://cnb.cool/medianexapp/plugin_api
    ref: main
gen:
  - plugin: go
    out: .
    opt: 
      - paths=source_relative
input:
  - local: testdata/proto/proto1.proto
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
	fields := strings.Fields(cmds[0])

	cmd := exec.Command(fields[0], fields[1:]...)
	_, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Stat("testdata/proto/proto1.pb")
	if err != nil {
		t.Fatal(err)
	}
	os.Remove("testdata/proto/proto1.pb")
	_, err = os.Stat("testdata/proto/proto1.pb.go")
	if err != nil {
		t.Fatal(err)
	}
	os.Remove("testdata/proto/proto1.pb.go")
}
