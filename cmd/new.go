package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new <title>",
	Short: "Create a new note",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Title required")
		}
		title := args[0]
		fmt.Printf("new called with %v\n", title)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
