package config

import (
	"os"

	"github.com/spf13/viper"

	"github.com/lukso-network/tools-lukso-cli/common/file"
)

type Configurator interface {
	Create(cfg NodeConfig) error
	Write(cfg NodeConfig) error
	Get() (cfg NodeConfig)
	Exists() bool
}

type config struct {
	path  string
	viper *viper.Viper
	file  file.Manager
	cfg   NodeConfig
}

var _ Configurator = &config{}

// NodeConfig represents the structure of the node folder configuration.
type NodeConfig struct {
	UseClients UseClients `mapstructure:"useclients"`
	Ipv4       string     `mapstructure:"ipv4"`
}

type UseClients struct {
	ExecutionClient string `mapstructure:"execution"`
	ConsensusClient string `mapstructure:"consensus"`
	ValidatorClient string `mapstructure:"validator"`
}

func NewConfigurator(path string) Configurator {
	dir, file, extension := parsePath(path)
	cfg := viper.New()

	cfg.AddConfigPath(dir)
	cfg.SetConfigName(file)
	cfg.SetConfigType(extension)

	return &config{
		path:  path,
		viper: cfg,
	}
}

// Create creates a new config that keeps track of selected dependencies and writes to it.
// By default, this file should be present in root of initialized lukso directory
func (c *config) Create(cfg NodeConfig) (err error) {
	err = c.file.Create(c.path)
	if err != nil {
		return
	}

	c.viper.Set("useClients.execution", cfg.UseClients.ExecutionClient)
	c.viper.Set("useClients.consensus", cfg.UseClients.ConsensusClient)
	c.viper.Set("useClients.validator", cfg.UseClients.ValidatorClient)
	c.viper.Set("ipv4", cfg.Ipv4)

	err = c.viper.WriteConfigAs(c.path)

	return
}

func (c *config) Exists() bool {
	_, err := os.Stat(c.path)

	return err == nil
}

func (c *config) Write(cfg NodeConfig) (err error) {
	err = c.viper.ReadInConfig()
	if err != nil {
		return
	}
	c.viper.Set("ipv4", "stub")

	err = c.viper.WriteConfigAs(c.path)

	return
}

// Read reads from config file passed during config instance into c
func (c *config) Read() (err error) {
	err = c.viper.ReadInConfig()
	if err != nil {
		return
	}

	err = c.viper.Unmarshal(&c.cfg)

	return
}

func (c *config) Get() NodeConfig {
	return c.cfg
}
