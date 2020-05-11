package config

type Config struct {
	SDK struct {
		MaxWaitBlock int64  `yaml:"max_wait_block"`
		RestEndpoint string `yaml:"rest_endpoint"`
		CliEndpoint  string `yaml:"cli_endpoint"`
	} `yaml:"sdk"`
	App struct {
	} `yaml:"app"`
}
