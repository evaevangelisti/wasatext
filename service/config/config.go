package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type WebAPIConfig struct {
	Config struct {
		FilePath string `conf:"default:./conf/config.yaml"`
	}

	Web struct {
		APIHost         string        `conf:"default:0.0.0.0:3000"`
		DebugHost       string        `conf:"default:0.0.0.0:4000"`
		ReadTimeout     time.Duration `conf:"default:5s"`
		WriteTimeout    time.Duration `conf:"default:5s"`
		ShutdownTimeout time.Duration `conf:"default:5s"`
	}

	Database struct {
		FilePath       string `conf:"default:./tmp/wasatext.db"`
		MigrationsPath string `conf:"default:./service/database/migrations"`
	}

	Debug bool
}

func LoadConfig() (WebAPIConfig, error) {
	var config WebAPIConfig

	if err := conf.Parse(os.Args[1:], "CFG", &config); err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			usage, err := conf.Usage("CFG", &config)

			if err != nil {
				return config, fmt.Errorf("generating config usage: %w", err)
			}

			logger := logrus.New()
			logger.Info(usage)

			return config, conf.ErrHelpWanted
		}

		return config, fmt.Errorf("parsing config: %w", err)
	}

	file, err := os.Open(config.Config.FilePath)

	if err != nil && !os.IsNotExist(err) {
		return config, fmt.Errorf("opening config file: %w", err)
	} else if err == nil {
		yamlData, err := io.ReadAll(file)

		if err != nil {
			return config, fmt.Errorf("reading config file: %w", err)
		}

		err = yaml.Unmarshal(yamlData, &config)

		if err != nil {
			return config, fmt.Errorf("unmarshalling config file: %w", err)
		}

		file.Close()
	}

	return config, nil
}
