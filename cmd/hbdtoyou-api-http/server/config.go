package server

import (
	"fmt"
	"hbdtoyou/pkg/environment"
	"os"
	"strings"
)

// configFileName is the template for config file name.
//
// Value of {BASEENV} will be changed in to the current
// environment.
const configFileName = "config.{BASEENV}.yaml"

// configFileLocation is the location of config files.
const configFileLocation = `/etc/hbdtoyou-api-http/`

// ConfigFilePaths contains list of paths to look for the
// config file. Paths are sorted in ascending order based on
// priority.
var configFilePaths = []string{
	configFileLocation + configFileName,           // for production and staging environment
	"files" + configFileLocation + configFileName, // for development environment
}

// getConfigFilePath checks file paths in the available file
// paths and returns the fist valid file path found.
func getConfigFilePath() (string, error) {
	for _, path := range configFilePaths {
		path = strings.Replace(path, "{BASEENV}", string(environment.GetServiceEnv()), -1)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		// found
		return path, nil
	}
	return "", fmt.Errorf("can't find valid config filepath")
}
