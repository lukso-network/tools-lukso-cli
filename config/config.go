package config

import (
	"github.com/spf13/viper"
	"os"
	"strings"
)

const Path = "./cli-config.yml"

// parsePath returns a path to file, excluding file, and the file name itself.
// Path parameter is used as a full file path.
// This parse is used to create a viper instance
// Example: ./path/to/file.txt => (./path/to, file.txt)
func parsePath(path string) (dir, fileName, extension string) {
	var lastIndex int

	segments := strings.Split(path, "/")

	if len(segments) != 0 {
		lastIndex = len(segments) - 1
		fullFile := segments[lastIndex]

		splittedFile := strings.Split(fullFile, ".")
		switch len(splittedFile) {
		case 0:
			break
		case 1:
			fileName = splittedFile[0]
		case 2:
			fileName = strings.Split(fullFile, ".")[0]
			extension = strings.Split(fullFile, ".")[1]
		}
	}

	switch len(segments) {
	case 0:
		return

	case 1:
		dir = "." // it means path is just a fileName - in current dir

		return

	default:
		dir = strings.Join(segments[:lastIndex], "/")

		return
	}
}

type Config struct {
	path            string
	viper           *viper.Viper
	executionClient string `mapstructure:"executionclient"`
	consensusClient string `mapstructure:"consensusclient"`
}

// NewConfig creates and initializes viper config instance - it doesn't load config, to load use c.Read().
func NewConfig(path string) *Config {
	dir, file, extension := parsePath(path)
	cfg := viper.New()

	cfg.AddConfigPath(dir)
	cfg.SetConfigName(file)
	cfg.SetConfigType(extension)

	return &Config{
		path:  path,
		viper: cfg,
	}
}

// Create creates a new config that keeps track of selected dependencies and writes to it.
// By default, this file should be present in root of initialized lukso directory
func (c *Config) Create(selectedExecution, selectedConsensus string) (err error) {
	_, err = os.Create(c.path)
	if err != nil {
		return
	}

	c.viper.Set("selectedclients.execution", selectedExecution)
	c.viper.Set("selectedclients.consensus", selectedConsensus)

	err = c.viper.WriteConfigAs(c.path)

	return
}

func (c *Config) Exists() bool {
	_, err := os.Stat(c.path)

	return err == nil
}

func (c *Config) WriteExecution(selectedExecution string) (err error) {
	c.viper.Set("selectedclients.execution", selectedExecution)

	err = c.viper.WriteConfigAs(c.path)

	return
}

func (c *Config) WriteConsensus(selectedConsensus string) (err error) {
	c.viper.Set("selectedclients.consensus", selectedConsensus)

	err = c.viper.WriteConfigAs(c.path)

	return
}

// Read reads from config file passed during config instance into c
func (c *Config) Read() (err error) {
	err = c.viper.ReadInConfig()
	if err != nil {
		return
	}

	c.executionClient = c.viper.Get("selectedclients.execution").(string)
	c.consensusClient = c.viper.Get("selectedclients.consensus").(string)

	return
}

func (c *Config) Execution() string {
	return c.executionClient
}

func (c *Config) Consensus() string {
	return c.consensusClient
}
