package config

import (
	"encoding/json"
	"os"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	kfile "github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"

	"github.com/lukso-network/tools-lukso-cli/common"
	"github.com/lukso-network/tools-lukso-cli/common/file"
)

type Configurator interface {
	create(cfg NodeConfig) error
	write() error
	get() (cfg NodeConfig)
	set(cfg NodeConfig) error
	exists() bool
}

type config struct {
	k            *koanf.Koanf
	fileProvider koanf.Provider
	parser       koanf.Parser

	path string
	file file.Manager

	cfg    NodeConfig
	valMap map[string]any
}

var (
	_            Configurator = &config{}
	configurator Configurator
)

// NodeConfig represents the structure of the node folder configuration.
// Even tho the config file is in YAML format, we use JSON tags for quick unmarshalling between types.
type NodeConfig struct {
	UseClients UseClients `json:"useclients"`
	Ipv4       string     `json:"ipv4"`
}

type UseClients struct {
	ExecutionClient string `json:"execution"`
	ConsensusClient string `json:"consensus"`
	ValidatorClient string `json:"validator"`
}

func UseConfigurator(c Configurator) {
	configurator = c
}

func NewConfigurator(path string, file file.Manager) Configurator {
	k := koanf.New(".")
	valMap := make(map[string]any)

	return &config{
		k:            k,
		fileProvider: kfile.Provider(path),
		valMap:       valMap,
		parser:       yaml.Parser(),
		path:         path,
		file:         file,
	}
}

func Create(cfg NodeConfig) error {
	return configurator.create(cfg)
}

func Write() error {
	return configurator.write()
}

func Get() (cfg NodeConfig) {
	return configurator.get()
}

func Set(cfg NodeConfig) error {
	return configurator.set(cfg)
}

func Exists() bool {
	return configurator.exists()
}

// Create creates a new config that keeps track of selected dependencies and writes to it.
// By default, this file should be present in root of initialized lukso directory
func (c *config) create(cfg NodeConfig) (err error) {
	err = c.file.Create(c.path)
	if err != nil {
		return
	}

	err = c.k.Load(c.fileProvider, c.parser)
	if err != nil {
		return
	}

	err = c.set(cfg)
	if err != nil {
		return
	}

	return c.write()
}

func (c *config) exists() bool {
	_, err := os.Stat(c.path)

	return err == nil
}

// Write writes the in-memory map to a file.
func (c *config) write() (err error) {
	parsed, err := c.k.Marshal(c.parser)
	if err != nil {
		return
	}

	return c.file.Write(c.path, parsed, common.ConfigPerms)
}

// Read reads from config file passed during config instance into c and returns the config
func (c *config) read() (err error) {
	err = c.k.Load(c.fileProvider, c.parser)
	if err != nil {
		return
	}

	b, err := json.Marshal(c.cfg)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &c.valMap)
	if err != nil {
		return
	}

	err = c.k.Load(confmap.Provider(c.valMap, "."), nil)
	if err != nil {
		return
	}

	err = c.k.UnmarshalWithConf("", &c.cfg, koanf.UnmarshalConf{Tag: "json"})

	return
}

// Get returns the in memory config.
func (c *config) get() NodeConfig {
	return c.cfg
}

// Set writes the config to the in memory state.
func (c *config) set(cfg NodeConfig) (err error) {
	b, err := json.Marshal(c.cfg)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &c.valMap)
	if err != nil {
		return
	}

	err = c.k.Load(confmap.Provider(c.valMap, "."), nil)
	if err != nil {
		return
	}

	c.cfg = cfg

	return
}
