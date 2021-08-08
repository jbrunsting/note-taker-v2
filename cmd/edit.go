package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jbrunsting/note-taker-v2/editor"
	"github.com/jbrunsting/note-taker-v2/html"
	"github.com/jbrunsting/note-taker-v2/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func writeHtml(dir string) error {
	files, err := utils.GetNotesFiles(dir)
	if err != nil {
		return err
	}
	htmlCode, err := html.GenerateHTML(files, dir)
	if err != nil {
		return err
	}
	htmlPath := fmt.Sprintf("%v/index.html", dir)
	_, err = os.Stat(htmlPath)
	if os.IsNotExist(err) {
		file, err := os.Create(htmlPath)
		if err != nil {
			return err
		}
		file.Close()
	}
	return ioutil.WriteFile(htmlPath, []byte(htmlCode), 0644)
}

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit an existing note",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := utils.GetDirFromCmd(cmd)
		if err != nil {
			return err
		}
		result, err := utils.SearchForFile(dir)
		if err != nil {
			return err
		}
		err = editor.Edit(result.Path)
		if err != nil {
			return err
		}
		return writeHtml(dir)
	},
}

func init() {
	editCmd.PersistentFlags().StringP("dir", "d", "", "Notes directory")
	viper.BindPFlag("dir", editCmd.PersistentFlags().Lookup("dir"))

	rootCmd.AddCommand(editCmd)
}
