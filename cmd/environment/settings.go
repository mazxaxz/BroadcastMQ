package environment

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Settings struct {
	ConfigPath  string `default:"/etc/broadcastmq/config.yaml"`
	LogLevel    string `default:"info"`
	OutputType  string `default:"text"`
}

func (s *Settings) LoadSettings() error {
	err := envconfig.Process("BMQ", s)
	if err != nil {
		return fmt.Errorf("An error occured while deserializing environment Settings: %v", err)
	}

	return nil
}

func (s *Settings) ConfigureLogging(log *logrus.Logger) {
	log.SetLevel(s.getLogLevel())
	if s.OutputType == "json" {
		log.SetFormatter(&logrus.JSONFormatter{})
	}
}

func (s *Settings) getLogLevel() logrus.Level {
	switch s.LogLevel {
	case "trace":
		return logrus.TraceLevel
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.InfoLevel
	}
}