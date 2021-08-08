package utils

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func GetDirFromCmd(cmd *cobra.Command) (string, error) {
	dir, err := cmd.Flags().GetString("dir")
	if err != nil {
		return "", err
	}
	if dir == "" {
		homedir, err := os.UserHomeDir()
		if err != nil {
			return "", &ReadableError{
				Err: err,
				Msg: "Could not find the home directory, and no notes directory provided - please provide the notes directory explicitly with the --dir flag\n",
			}
		}
		return fmt.Sprintf("%s/.note-taker-v2/notes", homedir), nil
	}
	return dir, nil
}

type NoteSearchResult struct {
	Name    string
	Path    string
	Preview string
}

func getSearchOptions(dir string) ([]NoteSearchResult, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return []NoteSearchResult{}, err
	}
	searchOptions := []NoteSearchResult{}
	for _, f := range files {
		filename := f.Name()
		if filename[len(filename)-7:] == "deleted" {
			continue
		}
		name := strings.TrimSuffix(filename, filepath.Ext(filename))
		path := fmt.Sprintf("%v/%v", dir, filename)
		file, err := os.Open(path)
		if err != nil {
			return []NoteSearchResult{}, err
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		preview := ""
		lines := 0
		for scanner.Scan() {
			preview = preview + scanner.Text() + "\n"
			lines += 1
			if lines >= 5 {
				preview = preview + "..."
				lines += 1
				break
			}
		}
		for ; lines < 5; lines += 1 {
			preview = preview + "\n"
		}
		searchOptions = append(searchOptions, NoteSearchResult{Name: name, Path: path, Preview: preview})
	}
	return searchOptions, nil
}

func SearchForFile(dir string) (*NoteSearchResult, error) {
	files, err := getSearchOptions(dir)
	if err != nil {
		return nil, err
	}
	templates := &promptui.SelectTemplates{
		Label:    "Notes:",
		Active:   "> {{ .Name | green }}",
		Inactive: "  {{ .Name | cyan }}",
		Selected: "> {{ .Name | green }}",
		Details: `
{{ .Preview | faint }}`,
	}

	searcher := func(input string, index int) bool {
		pepper := files[index]
		name := strings.Replace(strings.ToLower(pepper.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:             "Note",
		Items:             files,
		Templates:         templates,
		Size:              15,
		Searcher:          searcher,
		StartInSearchMode: true,
	}
	i, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	return &files[i], nil
}
