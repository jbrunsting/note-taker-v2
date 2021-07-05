package editor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const DefaultEditor = "vim"

func GetPath(dir string, title string) string {
	return fmt.Sprintf("%s/%s.md", strings.TrimRight(dir, "/"), title)
}

func Edit(dir string, title string) error {
	path := GetPath(dir, title)

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = DefaultEditor
	}

	executable, err := exec.LookPath(editor)
	if err != nil {
		return err
	}

	cmd := exec.Command(executable, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
