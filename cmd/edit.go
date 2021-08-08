package cmd

import (
	"github.com/jbrunsting/note-taker-v2/editor"
	"github.com/jbrunsting/note-taker-v2/utils"
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
		result, err := utils.SearchForFile(dir)
		if err != nil {
			return err
		}
		return editor.Edit(result.Path)
	},
}

func init() {
	editCmd.PersistentFlags().StringP("dir", "d", "", "Notes directory")
	viper.BindPFlag("dir", editCmd.PersistentFlags().Lookup("dir"))

	rootCmd.AddCommand(editCmd)
}
