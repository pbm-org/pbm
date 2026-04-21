package builder

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

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
		if dep.Remote == "" && dep.Local == "" {
			return fmt.Errorf("dep config is invalid")
		}
		if dep.Local != "" {
			stat, err := os.Stat(dep.Local)
			if err != nil {
				return err
			}
			if !stat.IsDir() {
				return fmt.Errorf("dep local must be dir")
			}
		}
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
		if input.Remote == "" && input.Local == "" {
			continue
		}
		depPaths = append(depPaths, input.PbPath)
	}

	for _, depPath := range depPaths {
		if depPath.Local != "" {
			_, err := os.Stat(depPath.Local)
			if err != nil {
				return err
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
			stat, err := os.Stat(depPath)
			if err != nil {
				return nil, err
			}
			if !stat.IsDir() {
				return nil, fmt.Errorf("dependency %s is not a directory", depPath)
			}
		}
		fmt.Fprintf(b, "-I %s ", depPath)
	}
	for _, input := range pbmCfg.Input {
		if input.Remote != "" {
			depPath := deps.GetDepDir(input.PbPath)
			fmt.Fprintf(b, "-I %s ", depPath)
		} else {
			depPath := input.Local
			_, err := os.Stat(depPath)
			if err != nil {
				return nil, err
			}
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
			stat, err := os.Stat(input.Local)
			if err != nil {
				return nil, err
			}
			if stat.IsDir() {
				protoFiles := []string{}
				filepath.Walk(input.Local, func(path string, info fs.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if info.IsDir() {
						return nil
					}
					if strings.HasSuffix(path, ".proto") {
						protoFiles = append(protoFiles, path)
					}

					return nil
				})
				depPath = strings.Join(protoFiles, " ")
			} else {
				if !strings.HasSuffix(depPath, ".proto") {
					return nil, fmt.Errorf("valid proto file")
				}
			}
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

func RunPbmCmd(cmds []string) error {
	g := errgroup.Group{}
	for _, cmd := range cmds {
		slog.Debug("run pbbuild cmd", "cmd", cmd)
		g.Go(func() error {
			fields := strings.Fields(cmd)
			buf := &bytes.Buffer{}
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			cmd := exec.CommandContext(ctx, fields[0], fields[1:]...)
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
