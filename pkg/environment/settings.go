package environment

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Settings struct {
	ConfigPath string `default:"/etc/broadcastmq/config.yaml"`
}

func (s *Settings) LoadSettings() error {
	err := envconfig.Process("BMQ", s)
	if err != nil {
		fmt.Errorf("An error occured while deserializing environment Settings: %v", err)
		return err
	}

	return nil
}