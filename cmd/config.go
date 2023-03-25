package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

// EndpointDefinition contains information about
// custom endpoints defined in YAML configuration files
type EndpointDefinition struct {
	Path         string `yaml:"path"`
	Method       string `yaml:"method"`
	Description  string `yaml:"description"`
	functionName string
}

// Returns a templated string describing a handler function
// for the given endpoint
func (e *EndpointDefinition) GenerateHandlerFunction() string {
	// Set the function name to a camelCase alphanumeric representation
	// of the endpoint path

	e.functionName = fmt.Sprintf("%s%sHandler", capitalizeAfterSlash(e.Path), e.Method)

	templateFunction := `
	// %s
	func (a *application) %s(w http.ResponseWriter, r *http.Request) {
		replyTextContent(w, r, http.StatusOK, "OK")
	} 
	`
	return fmt.Sprintf(templateFunction, e.Description, e.functionName)
}

// Returns a templated string to append to README.md API Endpoint table
func (e *EndpointDefinition) GenerateReadmeTableEntry() string {
	return fmt.Sprintf("| %s | %s | %s |\n", e.Path, e.Method, e.Description)
}

// Returns a templated line of code attaching handler function to path via router
func (e *EndpointDefinition) GenerateRouterAttachment() string {
	// ToDo: Format method into http.MethodXYZ form rather than just XYZ
	return fmt.Sprintf("router.HandlerFunc(%s, %s, a.%s)", e.Method, e.Path, e.functionName)
}

// YamlConfig contains app information collected
// from YAML configuration file
type YamlConfig struct {
	AppName   string               `yaml:"appName"`
	Directory string               `yaml:"directory"`
	ModName   string               `yaml:"modName"`
	Endpoints []EndpointDefinition `yaml:"endpoints"`
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

// Returns an alphanumeric camelcase representation of a
// path-style endpoint (i.e. /v1/healthcheck -> v1Healthcheck)
func capitalizeAfterSlash(s string) string {
	var result string
	capitalize := true

	for _, c := range s {
		if c == '/' {
			capitalize = true
		} else {
			if capitalize {
				result += strings.ToUpper(string(c))
				capitalize = false
			} else {
				result += string(c)
			}
		}
	}

	return result
}
