package utils

import (
	"fmt"
	"os"

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
