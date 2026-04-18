package deps

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pbm-org/pbm/internal/config"
)

func TestDeps(t *testing.T) {
	cacheDir, _ = os.Getwd()
	dep := config.Dep{
		PbPath: config.PbPath{
			Remote: "https://cnb.cool/medianexapp/plugin_api",
			Ref:    "main",
		},
	}

	err := CloneDepPath(dep.PbPath)
	if err != nil {
		t.Fatal(err)
	}
	depDir := GetDepDir(dep.PbPath)

	defer os.RemoveAll(filepath.Dir(filepath.Dir(depDir)))
	dirs, err := os.ReadDir(depDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(dirs) == 0 {
		t.Fatal("dir is empty")
	}
}
