package deps

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pbm-org/pbm/internal/config"
)

var (
	cacheDir, _ = os.UserCacheDir()
)

func GetCacheDir() (string, error) {
	var err error
	pdmCacheDir := filepath.Join(cacheDir, "pbm")
	_, err = os.Stat(pdmCacheDir)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}
		err = os.MkdirAll(pdmCacheDir, 0755)
		if err != nil {
			return "", err
		}
	}
	return pdmCacheDir, nil
}

func GetDepDir(dep config.PbPath) string {
	dir, err := GetCacheDir()
	if err != nil {
		panic(err)
	}
	depPath := filepath.Join(dir, url.QueryEscape(dep.Remote))
	if dep.Ref != "" {
		depPath = filepath.Join(depPath, dep.Ref)
	}
	return depPath
}

func CloneDepPath(dep config.PbPath) (err error) {
	depDir := GetDepDir(dep)
	dirs, err := os.ReadDir(depDir)

	if err == nil && len(dirs) > 0 {
		slog.Debug("dep already exist", "remote", dep.Remote, "ref", dep.Ref)
		return nil
	}
	err = os.MkdirAll(filepath.Dir(depDir), 0755)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			rePath := depDir
			if dep.Ref != "" {
				rePath = filepath.Dir(rePath)
			}
			os.RemoveAll(rePath)
		}
	}()
	cmdParams := []string{"git", "clone", "--depth", "1"}
	if dep.Ref != "" {
		cmdParams = append(cmdParams, "--branch", dep.Ref)
	}
	cmdParams = append(cmdParams, dep.Remote, depDir)
	slog.Debug("clone dep", "cmd", strings.Join(cmdParams, " "))
	err = runCmd(cmdParams)
	if err != nil {
		slog.Error("clone failed", "err", err)
		return err
	}
	slog.Debug("clone success", "path", depDir)
	return nil

}

func runCmd(cmdParams []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	cmd := exec.CommandContext(ctx, cmdParams[0], cmdParams[1:]...)
	buf := &bytes.Buffer{}
	cmd.Stderr = buf
	cmd.Stdout = buf
	defer cancel()
	cmd.Cancel = func() error {
		err := cmd.Process.Kill()
		return err
	}
	err := cmd.Run()
	if err != nil {
		fmt.Println("111")
		return fmt.Errorf("run cmd failed %s %s", err, buf.String())
	}
	fmt.Println("run over")
	return nil
}
