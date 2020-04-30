package configuration

import (
	"fmt"
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
	ConnectionString string `yaml:"connectionString"`
	Exchange         string `yaml:"exchange"`
	RoutingKey       string `yaml:"routingKey"`
}

type Destination struct {
	ConnectionString string  `yaml:"connectionString"`
	Queues           []Queue `yaml:"queues"`
}

type Queue struct {
	Name        string            `yaml:"name"`
	KeepHeaders bool              `yaml:"keepHeaders"`
	Args        map[string]string `yaml:"arguments"`
}

func (cfg *Config) LoadConfiguration(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("File at path: '%s' does not exist and will be created", path)
		mkdirErr := ensurePath(path, os.ModePerm)
		if mkdirErr != nil {
			return mkdirErr
		}
	}

	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Errorf("File with path: '%s' could not be opened. Error: %v", path, err)
		return err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Errorf("File with path: '%s' could not be read. Error: %v", path, err)
		return err
	}

	if len(b) == 0 {
		b = []byte("---")
		file.Write(b)
	}

	err = yaml.Unmarshal(b, cfg)
	if err != nil {
		fmt.Errorf("File with path: '%s' could parsed. Error: %v", path, err)
		return err
	}

	return nil
}

// default config values
func (cfg *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawConfig Config
	raw := rawConfig{
		LivenessProbe: Probe{
			Enabled: false,
			Path: "/health",
		},
		ReadinessProbe: Probe{
			Enabled: false,
			Path: "/ready",
		},
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
		fmt.Errorf("Could not create directory: '%s'. Error: %v", path, err)
		return err
	}

	return nil
}
