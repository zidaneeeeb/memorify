package config

import configlib "hbdtoyou/pkg/config"

type Payment struct {
	HTTP map[string]PaymentHTTP `yaml:"http"`
}

type PaymentHTTP struct {
	Timeout configlib.Duration `yaml:"timeout"`
}
