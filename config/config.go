package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config is a struct to manage game configuration
type Config struct {
	SDK struct {
		MaxWaitBlock int64  `yaml:"max_wait_block"`
		RestEndpoint string `yaml:"rest_endpoint"`
		CliEndpoint  string `yaml:"cli_endpoint"`
	} `yaml:"sdk"`
	App struct {
		DaemonTimeoutCommit int `yaml:"daemon_timeout_commit"`
	} `yaml:"app"`
	Terminal struct {
		UseLocalDm    bool
		UseRestTx     bool
		AutomateInput bool
	}
}

// ReadConfig is a function to read configuration
func ReadConfig() (Config, error) {
	args := os.Args

	var useRestTx bool = false
	var useLocalDm bool = false
	var automateInput bool = false

	if len(args) > 1 {
		for _, arg := range args[2:] {
			switch arg {
			case "-locald":
				useLocalDm = true
			case "-userest":
				useRestTx = true
			case "-automate":
				automateInput = true
			}
		}
	}

	cfgFileName := "config.yml"
	if useLocalDm {
		cfgFileName = "config_local.yml"
	}

	var cfg Config
	configf, err := os.Open(cfgFileName)
	if err == nil {
		defer configf.Close()

		decoder := yaml.NewDecoder(configf)
		err = decoder.Decode(&cfg)
		cfg.Terminal.UseLocalDm = useLocalDm
		cfg.Terminal.UseRestTx = useRestTx
		cfg.Terminal.AutomateInput = automateInput
	}
	return cfg, err
}
