package config

import configlib "hbdtoyou/pkg/config"

type Template struct {
	HTTP map[string]TemplateHTTP `yaml:"http"`
}

type TemplateHTTP struct {
	Timeout configlib.Duration `yaml:"timeout"`
}
