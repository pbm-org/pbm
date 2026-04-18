package deps

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

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
	_, err = os.Stat(depDir)
	if err == nil {
		return nil
	}
	err = os.MkdirAll(depDir, 0755)
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
	cmdParams := []string{"clone", "--depth", "1"}
	if dep.Ref != "" {
		cmdParams = append(cmdParams, "--branch", dep.Ref)
	}
	cmdParams = append(cmdParams, dep.Remote, depDir)
	slog.Info("start clone dep", "cmd", cmdParams)
	cmd := exec.Command("git", cmdParams...)
	buf := &bytes.Buffer{}
	cmd.Stderr = buf
	err = cmd.Run()
	if err != nil {
		cmd := exec.Command("git", "clone", dep.Remote, depDir)
		buf.Reset()
		cmd.Stderr = buf
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("clone failed %s %s", err, buf.String())
		}
		if dep.Ref != "" {
			err = os.Chdir(depDir)
			if err != nil {
				return err
			}
			cmd = exec.Command("git", "checkout", dep.Ref)
			cmd.Stderr = buf
			err = cmd.Run()
			if err != nil {
				return fmt.Errorf("clone failed %s %s", err, buf.String())
			}
		}
	}
	slog.Debug("clone %s %s success", dep.Remote, dep.Ref)
	return nil

}
