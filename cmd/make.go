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

func makeAction(out io.Writer, appName, modName, dir string) error {
	fmt.Printf("Creating new skeleton server named %s in %s\n", appName, dir)
	// Check given directory
	fmt.Printf("Checking if %s is a valid directory...\n", dir)
	_, err := os.Stat(dir)
	if err != nil {
		fmt.Printf("--> Couldn't find directory named %s, aborting.\n", dir)
		return err
	} else {
		fmt.Printf("--> Successfully confirmed %s exists, continuing.\n", dir)
	}
	// Check and create target directory
	target := filepath.Join(dir, appName)
	fmt.Printf("Creating application subdirectory %s...\n", target)
	err = os.Mkdir(target, 0755)
	if err != nil {
		fmt.Printf("--> Couldn't create directory %s, aborting.\n", target)
		return err
	} else {
		fmt.Printf("--> Successfully created directory %s, continuing\n", target)
	}
	// Create go mod
	fmt.Printf("Creating go module named %s in %s...\n", modName, target)
	cmd := exec.Command("go", "mod", "init", modName)
	cmd.Dir = target
	_, err = cmd.Output()
	if err != nil {
		fmt.Printf("--> Couldn't create go module %s, aborting.\n", modName)
		return err
	} else {
		fmt.Printf("--> Successfully created go module %s, continuing\n", modName)
	}
	// Create base files
	fmt.Println("Creating README.md")
	readme, err := os.Create(filepath.Join(target, "README.md"))
	if err != nil {
		fmt.Println("--> Couldn't create README.md, aborting.")
		return err
	} else {
		fmt.Println("--> Successfully created README.md, continuing")
	}
	defer readme.Close()
	_, err = readme.Write([]byte(fmt.Sprintf("# %s\n\n", appName)))
	if err != nil {
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
	_, err = readme.Write([]byte("## File Structure\n\n"))
	if err != nil {
		return err
	}
	for _, t := range targets {
		targ := filepath.Join(target, t[0])
		fmt.Printf("Creating subdirectory %s...\n", targ)
		err = os.Mkdir(targ, 0755)
		if err != nil {
			fmt.Printf("--> Couldn't create directory %s, aborting.\n", targ)
			return err
		} else {
			fmt.Printf("--> Successfully created directory %s, continuing\n", targ)
		}
		if t[1] != "" {
			_, err = readme.Write([]byte(fmt.Sprintf("- `%s`: %s\n", t[0], t[1])))
			if err != nil {
				return err
			}
		}
	}
	fmt.Printf("Creating %s/cmd/api/main.go...\n", target)
	main_dest, err := os.Create(filepath.Join(target, "cmd/api/main.go"))
	if err != nil {
		fmt.Println("--> Couldn't create cmd/api/main.go, aborting.")
		return err
	} else {
		fmt.Println("--> Successfully created cmd/api/main.go continuing")
	}
	main_source, err := os.Open("./skeleton-files/main.go")
	if err != nil {
		return err
	}
	_, err = io.Copy(main_dest, main_source)
	if err != nil {
		return err
	}
	defer main_source.Close()
	defer main_dest.Close()

	fmt.Printf("Creating %s/cmd/api/healthcheck.go...\n", target)
	hc_dest, err := os.Create(filepath.Join(target, "/cmd/api/healthcheck.go"))
	if err != nil {
		fmt.Println("--> Couldn't create /cmd/api/healthcheck.go, aborting.")
		return err
	} else {
		fmt.Println("--> Successfully created /cmd/api/healthcheck.go continuing")
	}
	hc_source, err := os.Open("./skeleton-files/healthcheck.go")
	if err != nil {
		return err
	}
	_, err = io.Copy(hc_dest, hc_source)
	if err != nil {
		return err
	}
	defer hc_source.Close()
	defer hc_dest.Close()

	fmt.Printf("Creating %s/Makefile...\n", target)
	m_dest, err := os.Create(filepath.Join(target, "Makefile"))
	if err != nil {
		fmt.Println("--> Couldn't create Makefile, aborting.")
		return err
	} else {
		fmt.Println("--> Successfully created Makefile continuing")
	}
	m_source, err := os.Open("./skeleton-files/Makefile")
	if err != nil {
		return err
	}
	_, err = io.Copy(m_dest, m_source)
	if err != nil {
		return err
	}
	defer main_source.Close()
	defer main_dest.Close()

	return nil
}

func init() {
	rootCmd.AddCommand(makeCmd)
	makeCmd.Flags().StringP("app-name", "n", "", "Name of application (Required)")
	makeCmd.Flags().StringP("mod-name", "m", "", "Name of top-level application go module (default $app-name)")
	makeCmd.Flags().StringP("dir", "d", "./", "Path to target app directory")
}
