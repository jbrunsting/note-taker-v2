package cmd

import (
	"time"

	"github.com/jbrunsting/note-taker-v2/editor"
	"github.com/jbrunsting/note-taker-v2/html"
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
		writing := make(chan bool, 1)
		go func() {
			for {
				time.Sleep(5 * time.Second)
				writing <- true
				html.WriteHtml(dir)
				<-writing
			}
		}()
		err = editor.Edit(result.Path)
		if err != nil {
			return err
		}
		writing <- true
		return html.WriteHtml(dir)
	},
}

func init() {
	editCmd.PersistentFlags().StringP("dir", "d", "", "Notes directory")
	viper.BindPFlag("dir", editCmd.PersistentFlags().Lookup("dir"))

	rootCmd.AddCommand(editCmd)
}
