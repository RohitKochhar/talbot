package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

// Creates an empty file in a specified path
// Returns the file for future use if needed
func createFile(name, path string) (*os.File, error) {
	fp := filepath.Join(path, name)
	fmt.Printf("Creating file: %s...\n", fp)
	f, err := os.Create(fp)
	if err != nil {
		fmt.Printf("--> Couldn't create %s, aborting.\n", fp)
		return nil, err
	} else {
		fmt.Printf("--> Successfully created %s, continuing\n", fp)
	}
	return f, nil
}

// Writes a string to a given file
func writeFile(file *os.File, content string) error {
	_, err := file.Write([]byte(content + "\n"))
	return err
}
