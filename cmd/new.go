package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jbrunsting/note-taker-v2/editor"
	"github.com/jbrunsting/note-taker-v2/utils"
)

var newCmd = &cobra.Command{
	Use:   "new <title>",
	Short: "Create a new note",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Title argument required")
		}
		title := args[0]
		dir, err := utils.GetDirFromCmd(cmd)
		if err != nil {
			return err
		}
		filetype, err := cmd.Flags().GetString("filetype")
		if err != nil {
			return err
		}
		return createAndEdit(dir, title, filetype)
	},
}

func init() {
	newCmd.PersistentFlags().StringP("dir", "d", "", "Notes directory")
	newCmd.PersistentFlags().StringP("filetype", "f", "md", "Filetype to use for the note")
	viper.BindPFlag("dir", newCmd.PersistentFlags().Lookup("dir"))
	viper.BindPFlag("filetype", newCmd.PersistentFlags().Lookup("filetype"))

	rootCmd.AddCommand(newCmd)
}

func createAndEdit(dir string, title string, filetype string) error {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	if title == "index" && filetype == "html" {
		return nil
	}
	filepath := fmt.Sprintf("%s/%s.%s", strings.TrimRight(dir, "/"), title, filetype)
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		file, err := os.Create(filepath)
		if err != nil {
			return err
		}
		defer file.Close()
	} else {
		return &utils.ReadableError{
			Err: nil,
			Msg: "A note with that title already exists",
		}
	}
	return editor.Edit(filepath)
}
