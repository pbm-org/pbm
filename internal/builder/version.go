package builder

import (
	"bytes"
	"os/exec"
	"strings"
)

func checkProtoc() string {
	cmd := exec.Command("protoc", "--version")
	buf := bytes.Buffer{}
	cmd.Stdout = &buf
	err := cmd.Run()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(buf.String())
}
