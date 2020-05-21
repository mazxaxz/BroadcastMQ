package config

import (
	"fmt"
	"github.com/guiferpa/gody/v2"
	"github.com/guiferpa/gody/v2/rule"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

type Config struct {
	LivenessProbe  Probe       `yaml:"livenessProbe"`
	ReadinessProbe Probe       `yaml:"readinessProbe"`
	Broadcasts     []Broadcast `yaml:"broadcasts"`
}

type Probe struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
}

type Broadcast struct {
	Source      Source      `yaml:"source"`
	Destination Destination `yaml:"destination"`
}

type Source struct {
	ConnectionString string `yaml:"connectionString" validate:"not_empty"`
	Exchange         string `yaml:"exchange" validate:"not_empty"`
	RoutingKey       string `yaml:"routingKey" validate:"not_empty"`
	BMQQueueName     string `yaml:"bmqQueueName"`
}

type Destination struct {
	ConnectionString string  `yaml:"connectionString" validate:"not_empty"`
	BMQExchange      string  `yaml:"bmqExcahnge"`
	BMQRoutingKey    string  `yaml:"bmqRoutingKey"`
	Queues           []Queue `yaml:"queues"`
	PersistHeaders   bool    `yaml:"persistHeaders"`
}

type Queue struct {
	Name          string `yaml:"name" validate:"not_empty"`
	BMQBindingKey string `yaml:"bmqBindingKey"`
	EnsureExists  bool   `yaml:"ensureExists"`
}

// LoadConfiguration loads variables from YAML file
func (cfg *Config) LoadConfiguration(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		mkdirErr := ensurePath(path, os.ModePerm)
		if mkdirErr != nil {
			return mkdirErr
		}
	}

	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return fmt.Errorf("File with path: '%s' could not be opened. Error: %v", path, err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("File with path: '%s' could not be read. Error: %v", path, err)
	}

	if len(b) == 0 {
		b = []byte("---")
		file.Write(b)
	}

	err = yaml.Unmarshal(b, cfg)
	if err != nil {
		return fmt.Errorf("File with path: '%s' could parsed. Error: %v", path, err)
	}

	return nil
}

// Validate makes sure that configuration is valid
func (cfg *Config) Validate() error {
	validator := gody.NewValidator()

	rules := []gody.Rule{rule.NotEmpty}
	validator.AddRules(rules...)

	if _, err := validator.Validate(*cfg); err != nil {
		return err
	}

	return nil
}

// FillDefault makes sure that optional values are filled
func (cfg *Config) FillDefault() {
	for i := range cfg.Broadcasts {
		if cfg.Broadcasts[i].Source.BMQQueueName == "" {
			cfg.Broadcasts[i].Source.BMQQueueName = DefaultBMQQueue
		}

		if cfg.Broadcasts[i].Destination.BMQRoutingKey == "" {
			cfg.Broadcasts[i].Destination.BMQRoutingKey = DefaultBMQRoutingKey
		}

		if cfg.Broadcasts[i].Destination.BMQExchange == "" {
			cfg.Broadcasts[i].Destination.BMQExchange = DefaultBMQExchange
		}

		for qi := range cfg.Broadcasts[i].Destination.Queues {
			if cfg.Broadcasts[i].Destination.Queues[qi].BMQBindingKey == "" {
				cfg.Broadcasts[i].Destination.Queues[qi].BMQBindingKey = DefaultBMQBindingKey
			}
		}
	}
}

func (cfg *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawConfig Config
	raw := rawConfig{
		LivenessProbe: Probe{
			Enabled: false,
			Path:    "/health",
		},
		ReadinessProbe: Probe{
			Enabled: false,
			Path:    "/ready",
		},
		Broadcasts: make([]Broadcast, 0),
	}

	if err := unmarshal(&raw); err != nil {
		return err
	}

	*cfg = Config(raw)
	return nil
}

func ensurePath(path string, perm os.FileMode) error {
	separator := "/"
	if runtime.GOOS == "windows" {
		separator = "\\"
	}
	dirs := strings.Split(path, separator)
	if len(dirs) > 0 {
		dirs = dirs[:len(dirs)-1]
	}
	if len(dirs) == 0 {
		return nil
	}

	err := os.MkdirAll(strings.Join(dirs, separator), perm)
	if err != nil {
		return fmt.Errorf("Could not create directory: '%s'. Error: %v", path, err)
	}

	return nil
}
