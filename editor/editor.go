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

func Delete(dir string, title string) error {
	path := GetPath(dir, title)
	renamed := fmt.Sprintf("%v.deleted", path)
	i := 1
	for {
		_, err := os.Stat(renamed)
		if os.IsNotExist(err) {
			break
		} else if err != nil {
			return err
		}
		renamed = fmt.Sprintf("%v.%v.deleted", path, i)
		i = i + 1
	}
	return os.Rename(path, renamed)
}
