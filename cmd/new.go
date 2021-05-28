package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var newCmd = &cobra.Command{
	Use:   "new <title>",
	Short: "Create a new note",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Title required")
		}
		title := args[0]
		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			return err
		}

		fmt.Printf("new called with %v/%v\n", dir, title)
		return nil
	},
}

func init() {
	newCmd.PersistentFlags().StringP("dir", "d", "", "Notes directory")
	viper.BindPFlag("dir", newCmd.PersistentFlags().Lookup("dir"))
	cobra.MarkFlagRequired(newCmd.PersistentFlags(), "dir")

	rootCmd.AddCommand(newCmd)
}
