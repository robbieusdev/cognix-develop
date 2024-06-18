package utils

import "github.com/caarlos0/env/v10"

func ReadConfig(cfg interface{}) error {
	return env.Parse(cfg)
}
