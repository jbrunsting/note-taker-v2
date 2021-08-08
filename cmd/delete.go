package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jbrunsting/note-taker-v2/editor"
	"github.com/jbrunsting/note-taker-v2/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Soft delete a note",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := utils.GetDirFromCmd(cmd)
		if err != nil {
			return err
		}
		result, err := utils.SearchForFile(dir)
		if err != nil {
			return err
		}
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Printf("Delete '%s'? [y/n]: ", result.Name)

			response, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}

			response = strings.ToLower(strings.TrimSpace(response))

			if response == "y" || response == "yes" {
				return editor.Delete(result.Path)
			} else if response == "n" || response == "no" {
				fmt.Printf("Negative confirmation, canceling\n")
				return nil
			}
		}
	},
}

func init() {
	deleteCmd.PersistentFlags().StringP("dir", "d", "", "Notes directory")
	viper.BindPFlag("dir", deleteCmd.PersistentFlags().Lookup("dir"))

	rootCmd.AddCommand(deleteCmd)
}
