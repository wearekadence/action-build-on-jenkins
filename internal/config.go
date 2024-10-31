package internal

import (
	"github.com/caarlos0/env"
)

type Config struct {
	JenkinsURL   string `env:"INPUT_URL"`
	JobName      string `env:"INPUT_JOB_NAME"`
	Username     string `env:"INPUT_USERNAME"`
	ApiToken     string `env:"INPUT_API_TOKEN"`
	Parameters   string `env:"INPUT_PARAMETERS" envDefault:"{}"`
	Wait         bool   `env:"INPUT_WAIT"  envDefault:"true"`
	Timeout      int    `env:"INPUT_TIMEOUT" envDefault:"600"`
	StartTimeout int    `env:"INPUT_START_TIMEOUT" envDefault:"600"`
	Interval     int    `env:"INPUT_INTERVAL" envDefault:"5"`
}

func NewConfig() *Config {
	cfg := Config{}
	env.Parse(&cfg)
	return &cfg
}
