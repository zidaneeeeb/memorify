package config

import configlib "hbdtoyou/pkg/config"

type Content struct {
	HTTP map[string]ContentHTTP `yaml:"http"`
}

type ContentHTTP struct {
	Timeout configlib.Duration `yaml:"timeout"`
}
