package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Config interface must contain app information
// collected from command flags or YAML file
type Config interface {
	getAppName() string   // Returns name of application
	getDirectory() string // Returns target directory
	getModName() string   // Returns name of go module
}

// FlagConfig contains app information collected
// from command line flags
type FlagConfig struct {
	AppName   string
	Directory string
	ModName   string
}

// Returns name of application
func (c FlagConfig) getAppName() string {
	return c.AppName
}

// Returns target directory
func (c FlagConfig) getDirectory() string {
	return c.Directory
}

// Returns name of go module
func (c FlagConfig) getModName() string {
	return c.ModName
}

// YamlConfig contains app information collected
// from YAML configuration file
type YamlConfig struct {
	AppName   string `yaml:"appName"`
	Directory string `yaml:"directory"`
	ModName   string `yaml:"modName"`
}

// Returns name of application
func (c YamlConfig) getAppName() string {
	return c.AppName
}

// Returns target directory
func (c YamlConfig) getDirectory() string {
	return c.Directory
}

// Returns name of go module
func (c YamlConfig) getModName() string {
	return c.ModName
}

// Checks if the config should be loaded from a YAML
// or from command flags and returns the appropriate
// Config interface or an error if applicable
func setConfiguration(cmd *cobra.Command) (Config, error) {
	// Check if a YAML was provided
	confFile, err := cmd.Flags().GetString("config")
	if err != nil {
		return nil, err
	}
	if confFile != "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		return loadYamlConfig(filepath.Join(wd, confFile))
	}
	// Check if flags were provided
	return loadFlagConfig(cmd)

}

// Reads configuration information from YAML file and returns
// a YamlConfig object containing application information
func loadYamlConfig(filename string) (*YamlConfig, error) {
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	yamlConf := &YamlConfig{}
	err = yaml.Unmarshal(yamlFile, yamlConf)
	if err != nil {
		return nil, err
	}
	fmt.Printf("HERE: %+v", yamlConf)
	if yamlConf.AppName == "" {
		return nil, fmt.Errorf("no `appName` field set in YAML file")
	}
	if yamlConf.Directory == "" {
		yamlConf.Directory = "./"
	}
	if yamlConf.ModName == "" {
		yamlConf.ModName = yamlConf.AppName
	}
	return yamlConf, nil
}

// Reads configuration information from command line flags and
// returns a FlagConfig object containing application information
func loadFlagConfig(cmd *cobra.Command) (*FlagConfig, error) {
	appName, err := cmd.Flags().GetString("app-name")
	if err != nil {
		return nil, err
	}
	if appName == "" {
		return nil, fmt.Errorf("no app-name argument was provided")
	}
	modName, err := cmd.Flags().GetString("mod-name")
	if err != nil {
		return nil, err
	}
	if modName == "" {
		modName = appName
	}
	dir, err := cmd.Flags().GetString("dir")
	if err != nil {
		return nil, err
	}
	return &FlagConfig{
		AppName:   appName,
		ModName:   modName,
		Directory: dir,
	}, nil
}
