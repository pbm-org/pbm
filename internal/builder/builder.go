package builder

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pbm-org/pbm/internal/config"
	"github.com/pbm-org/pbm/internal/deps"
	"golang.org/x/sync/errgroup"
)

func CleabPbmDep(pbmCfg *config.PbmConfig) error {
	for _, dep := range pbmCfg.Deps {
		if dep.Remote != "" {
			err := os.RemoveAll(deps.GetDepDir(dep.PbPath))
			if err != nil {
				return err
			}
		}
	}

	for _, input := range pbmCfg.Input {
		if input.Remote != "" {
			err := os.RemoveAll(deps.GetDepDir(input.PbPath))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CheckPbCfg(pbmCfg *config.PbmConfig) error {
	_, err := exec.LookPath("protoc")
	if err != nil {
		return err
	}
	depPaths := []config.PbPath{}
	for _, dep := range pbmCfg.Deps {
		depPaths = append(depPaths, dep.PbPath)
	}

	for _, input := range pbmCfg.Input {
		_, err := os.Stat(filepath.Dir(input.DescOut))
		if err != nil {
			err = os.MkdirAll(filepath.Dir(input.DescOut), 0755)
			if err != nil {
				return err
			}
		}
		depPaths = append(depPaths, input.PbPath)
	}

	for _, depPath := range depPaths {
		if depPath.Local != "" {
			match, err := filepath.Glob(depPath.Local)
			if err != nil {
				return err
			}
			if len(match) == 0 {
				return fmt.Errorf("%s not match proto file", depPath.Local)
			}
		} else if depPath.Remote != "" {
			err := deps.CloneDepPath(depPath)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("%v config is valid", depPath)
		}
	}
	for _, gen := range pbmCfg.Gen {
		_, err := exec.LookPath(fmt.Sprintf("protoc-gen-%s", gen.Plugin))
		if err != nil {
			return err
		}
	}
	return nil
}

func PbBuildCmd(pbmCfg *config.PbmConfig) ([]string, error) {
	b := &bytes.Buffer{}
	fmt.Fprintf(b, "protoc ")
	for _, dep := range pbmCfg.Deps {
		depPath := ""
		if dep.Remote != "" {
			depPath = deps.GetDepDir(dep.PbPath)
		} else {
			depPath = dep.Local
		}
		fmt.Fprintf(b, "-I %s ", depPath)
		// fmt.Fprintf(b, " \\ \n")
	}
	for _, input := range pbmCfg.Input {
		if input.Remote != "" {
			depPath := deps.GetDepDir(input.PbPath)
			fmt.Fprintf(b, "-I %s ", depPath)
		} else {
			fmt.Fprintf(b, "-I . ")
		}
	}
	for _, gen := range pbmCfg.Gen {
		fmt.Fprintf(b, "--%s_out=%s ", gen.Plugin, gen.Out)
		if len(gen.Opt) > 0 {
			for _, opt := range gen.Opt {
				fmt.Fprintf(b, "--%s_opt=%s ", gen.Plugin, opt)
			}
		}
	}

	cmds := []string{}
	for _, input := range pbmCfg.Input {
		cmd := ""
		depPath := ""
		if input.Remote != "" {
			depPath = deps.GetDepDir(input.PbPath)
			depPath = filepath.Join(depPath, input.File)
		} else {
			depPath = input.Local
		}
		if input.DescOut != "" {
			cmd = fmt.Sprintf(" --descriptor_set_out=%s --include_imports --include_source_info", input.DescOut)
		}
		cmd += " " + depPath
		slog.Debug("build proto", "file", depPath)
		cmds = append(cmds, b.String()+cmd)
	}
	return cmds, nil
}

func PbmCmd(cmds []string) error {
	g := errgroup.Group{}
	for _, cmd := range cmds {
		slog.Debug("run pbbuild cmd", "cmd", cmd)
		g.Go(func() error {
			fields := strings.Fields(cmd)
			buf := &bytes.Buffer{}
			cmd := exec.Command(fields[0], fields[1:]...)
			cmd.Stdout = buf
			cmd.Stderr = buf
			err := cmd.Run()
			if err != nil {
				return errors.New(buf.String())
			}
			return nil
		})
	}
	err := g.Wait()
	return err
}
