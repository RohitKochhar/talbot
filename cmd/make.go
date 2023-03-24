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

	"github.com/spf13/cobra"
)

// makeCmd represents the make command
var makeCmd = &cobra.Command{
	Use:     "make",
	Aliases: []string{"m"},
	Short:   "Makes a new server with the given specifications",
	Long:    "Makes a new server with the given specifications",
	RunE: func(cmd *cobra.Command, args []string) error {
		appName, err := cmd.Flags().GetString("app-name")
		if err != nil {
			return err
		}
		if appName == "" {
			return fmt.Errorf("no app-name argument was provided")
		}
		modName, err := cmd.Flags().GetString("mod-name")
		if err != nil {
			return err
		}
		if modName == "" {
			modName = appName
		}
		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			return err
		}
		return makeAction(os.Stdout, appName, modName, dir)
	},
}

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

// Creates an empty file in a specified path
// Returns the file for future use if needed
func createFile(name, path string) (*os.File, error) {
	fp := filepath.Join(path, name)
	fmt.Printf("Creating file: %s...\n", fp)
	f, err := os.Create(fp)
	if err != nil {
		fmt.Printf("--> Couldn't create %s, aborting.", fp)
		return nil, err
	} else {
		fmt.Printf("--> Successfully created %s, continuing", fp)
	}
	return f, nil
}

// Writes a string to a given file
func writeFile(file *os.File, content string) error {
	_, err := file.Write([]byte(content + "\n"))
	return err
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

func makeAction(out io.Writer, appName, modName, dir string) error {
	fmt.Printf("Creating new skeleton server named %s in %s\n", appName, dir)
	// Check given directory
	if err := checkDirectory(dir); err != nil {
		return err
	}
	// Check and create target directory
	target := filepath.Join(dir, appName)
	if err := createTargetDirectory(target); err != nil {
		return err
	}
	// Create go mod
	if err := initializeGoMod(modName, target); err != nil {
		return err
	}
	// Create README
	readme, err := createFile("README.md", target)
	if err != nil {
		return err
	}
	if err := writeFile(readme, fmt.Sprintf("# %s\n\n", appName)); err != nil {
		return err
	}

	targets := [][]string{
		{"bin", "Contains compiled application binaries ready for production deployment"},
		{"cmd", ""},
		{"cmd/api", "Contains application specific code to run server"},
		{"internal", "Contains various ancillary packages used by API"},
		{"migrations", "Contains SQL migration files for database"},
		{"remote", "Contains configuration files and setup scripts for remote deployment"},
	}

	if err := ScaffoldProject(target, targets, readme); err != nil {
		return err
	}

	if err := TemplateFile(target, "cmd/api/main.go", "./skeleton-files/main.go"); err != nil {
		return err
	}

	if err := TemplateFile(target, "cmd/api/handlers.go", "./skeleton-files/handlers.go"); err != nil {
		return err
	}

	if err := TemplateFile(target, "cmd/api/handlers_test.go", "./skeleton-files/handlers_test.go"); err != nil {
		return err
	}

	if err := TemplateFile(target, "Dockerfile", "./skeleton-files/Dockerfile"); err != nil {
		return err
	}

	if err := TemplateFile(target, "docker-compose.yaml", "./skeleton-files/docker-compose.yaml"); err != nil {
		return err
	}

	if err := TemplateFile(filepath.Join(target, "/cmd/api"), "docker-compose.yaml", "./skeleton-files/docker-compose.yaml"); err != nil {
		return err
	}

	if err := GetGolangPackage(target, "github.com/julienschmidt/httprouter"); err != nil {
		return err
	}

	// Update readme with healthcheck API info
	if err := writeFile(readme, "\n## API Endpoints\n"); err != nil {
		return err
	}
	if err := writeFile(readme, "| HTTP Endpoint | Method | Info |\n|-----|------|------|\n|`/v1/healthcheck`| GET | Displays server status\n"); err != nil {
		return err
	}

	if err := TemplateFile(target, "Makefile", "./skeleton-files/Makefile"); err != nil {
		return err
	}

	// Write a note about auto-generated documentation
	if err := writeFile(readme, "\n## `skele-server` disclaimer\n"); err != nil {
		return err
	}
	if err := writeFile(readme, "This README has been autogenerated by [skele-server](https://github.com/rohitkochhar/skele-server)"); err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(makeCmd)
	makeCmd.Flags().StringP("app-name", "n", "", "Name of application (Required)")
	makeCmd.Flags().StringP("mod-name", "m", "", "Name of top-level application go module (default $app-name)")
	makeCmd.Flags().StringP("dir", "d", "./", "Path to target app directory")
}
