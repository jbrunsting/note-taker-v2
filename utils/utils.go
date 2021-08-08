package utils

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

const (
	resultsHeight = 7
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

type byModTime []os.FileInfo

func (a byModTime) Len() int           { return len(a) }
func (a byModTime) Less(i, j int) bool { return a[i].ModTime().After(a[j].ModTime()) }
func (a byModTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func GetNotesFiles(dir string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return []os.FileInfo{}, err
	}
	sort.Sort(byModTime(files))
	filteredFiles := []os.FileInfo{}
	for _, f := range files {
		filename := f.Name()
		if filename[len(filename)-7:] == "deleted" || filename == "index.html" {
			continue
		}
		filteredFiles = append(filteredFiles, f)
	}
	return filteredFiles, nil
}

func getSearchOptions(dir string) ([]NoteSearchResult, error) {
	previewWidth := -1
	previewHeight := -1
	if term.IsTerminal(0) {
		width, height, err := term.GetSize(0)
		if err == nil {
			previewWidth = width - 4 // 2 char border on each side
			// 6 = 2 char search/notes title, 1 char space, 2 char top/bottom border, 1 char bottom space
			previewHeight = height - resultsHeight - 6
		}
	}
	files, err := GetNotesFiles(dir)
	if err != nil {
		return []NoteSearchResult{}, err
	}
	searchOptions := []NoteSearchResult{}
	for _, f := range files {
		filename := f.Name()
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
		if previewWidth > 0 && previewHeight > 0 {
			for scanner.Scan() {
				line := scanner.Text()
				if len(line) > previewWidth-3 {
					line = strings.TrimSpace(line[:previewWidth-3]) + "..."
				}
				preview = preview + "| " + line + strings.Repeat(" ", previewWidth-len(line)) + " |\n"
				lines += 1
				if lines >= previewHeight {
					break
				}
			}
			for ; lines < previewHeight; lines += 1 {
				preview = preview + "| " + strings.Repeat(" ", previewWidth) + " |\n"
			}
			preview = "+" + strings.Repeat("-", previewWidth+2) + "+\n" + preview
			preview = preview + "+" + strings.Repeat("-", previewWidth+2) + "+"
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
		Size:              resultsHeight,
		Searcher:          searcher,
		StartInSearchMode: true,
	}
	i, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	return &files[i], nil
}
