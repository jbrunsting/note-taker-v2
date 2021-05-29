package editor

import (
	"fmt"
)

func Edit(dir string, title string) error {
	fmt.Printf("Editing %v %v\n", dir, title)
	return nil
}
