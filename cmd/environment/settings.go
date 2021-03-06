package environment

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Settings struct {
	ConfigPath string `default:"/etc/broadcastmq/config.yaml"`
	LogLevel   string `default:"info"`
	OutputType string `default:"text"`
}

// LoadSettings loads environment variables into application runtime
func (s *Settings) LoadSettings() error {
	if err := envconfig.Process("BMQ", s); err != nil {
		return errors.Wrap(err, "An error occured while deserializing environment Settings")
	}

	return nil
}

// EncodeFormatter encodes formatter type from environment string
func EncodeFormatter(outputType string) logrus.Formatter {
	if outputType == "json" {
		return &logrus.JSONFormatter{}
	}

	return &logrus.TextFormatter{}
}

// EncodeLogLevel encodes log level from environment string
func EncodeLogLevel(level string) logrus.Level {
	switch level {
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
