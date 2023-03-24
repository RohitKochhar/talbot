/*
Copyright © 2023 Rohit Singh rkochhar@uwaterloo.ca
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// Checks to see if the given directory exists
func checkDirectory(dir string) error {
	fmt.Printf("Checking if %s is a valid directory...\n", dir)
	_, err := os.Stat(dir)
	if err != nil {
		fmt.Printf("--> Couldn't find directory named %s, aborting.\n", dir)
		return err
	} else {
		fmt.Printf("--> Successfully confirmed %s exists, continuing.\n", dir)
	}
	return nil
}

// Creates target directory
func createTargetDirectory(target string) error {
	fmt.Printf("Creating application subdirectory %s...\n", target)
	err := os.Mkdir(target, 0755)
	if err != nil {
		fmt.Printf("--> Couldn't create directory %s, aborting.\n", target)
		return err
	} else {
		fmt.Printf("--> Successfully created directory %s, continuing\n", target)
	}
	return nil
}

// Initializes go module
func initializeGoMod(modName, target string) error {
	fmt.Printf("Creating go module named %s in %s...\n", modName, target)
	cmd := exec.Command("go", "mod", "init", modName)
	cmd.Dir = target
	_, err := cmd.Output()
	if err != nil {
		fmt.Printf("--> Couldn't create go module %s, aborting.\n", modName)
		return err
	} else {
		fmt.Printf("--> Successfully created go module %s, continuing\n", modName)
	}
	return nil
}

// Scaffolds the project file structure and templates README documentation
func ScaffoldProject(target string, folders [][]string, readme *os.File) error {
	if err := writeFile(readme, "## File Structure\n\n"); err != nil {
		return err
	}

	for _, f := range folders {
		t := filepath.Join(target, f[0])
		fmt.Printf("Creating subdirectory %s...\n", t)
		if err := os.Mkdir(t, 0755); err != nil {
			fmt.Printf("--> Couldn't create directory %s, aborting.\n", t)
			return err
		} else {
			fmt.Printf("--> Successfully created directory %s, continuing\n", t)
		}
		if f[1] != "" {
			if err := writeFile(readme, fmt.Sprintf("- `%s`: %s\n", f[0], f[1])); err != nil {
				return err
			}
		}
	}

	return nil
}

// Templates a new file from an existing one
func TemplateFile(target string, newFile string, templateFile string) error {
	newTarget := filepath.Join(target, newFile)
	fmt.Printf("Creating %s\n", newTarget)
	main_dest, err := os.Create(newTarget)
	if err != nil {
		fmt.Printf("--> Couldn't create %s, aborting.\n", newTarget)
		return err
	} else {
		fmt.Printf("--> Successfully created %s continuing\n", newTarget)
	}
	main_source, err := os.Open(templateFile)
	if err != nil {
		return err
	}
	_, err = io.Copy(main_dest, main_source)
	if err != nil {
		return err
	}
	defer main_source.Close()
	defer main_dest.Close()
	return nil
}

// Executes `go mod download`
func GetGolangPackage(target string, packageName string) error {
	fmt.Printf("Collecting go package %s in %s...\n", packageName, target)
	cmd := exec.Command("go", "get", packageName)
	cmd.Dir = target
	_, err := cmd.Output()
	if err != nil {
		fmt.Printf("--> Couldn't collect go package %s, aborting.\n", packageName)
		return err
	} else {
		fmt.Printf("--> Successfully collected go package %s, continuing.\n", packageName)
	}
	return nil
}
