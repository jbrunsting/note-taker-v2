package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/jbrunsting/note-taker-v2/editor"
	"github.com/jbrunsting/note-taker-v2/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit an existing note",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := utils.GetDirFromCmd(cmd)
		if err != nil {
			return err
		}
		return searchAndEdit(dir)
	},
}

func init() {
	editCmd.PersistentFlags().StringP("dir", "d", "", "Notes directory")
	viper.BindPFlag("dir", editCmd.PersistentFlags().Lookup("dir"))

	rootCmd.AddCommand(editCmd)
}

type searchOption struct {
	Name    string
	Preview string
}

func getSearchOptions(dir string) ([]searchOption, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return []searchOption{}, err
	}
	searchOptions := []searchOption{}
	for _, f := range files {
		filename := f.Name()
		name := strings.TrimSuffix(filename, filepath.Ext(filename))
		path := fmt.Sprintf("%v/%v", dir, filename)
		file, err := os.Open(path)
		if err != nil {
			return []searchOption{}, err
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
		searchOptions = append(searchOptions, searchOption{Name: name, Preview: preview})
	}
	return searchOptions, nil
}

func searchAndEdit(dir string) error {
	files, err := getSearchOptions(dir)
	if err != nil {
		return err
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
		return err
	}

	return editor.Edit(dir, files[i].Name)
}
