package config

import configlib "hbdtoyou/pkg/config"

type User struct {
	PasswordSalt    string              `yaml:"password_salt"`
	TokenExpiration configlib.Duration  `yaml:"token_expiration"`
	TokenSecretKey  string              `yaml:"token_secret_key"`
	ClientID        string              `yaml:"client_id"`
	HTTP            map[string]UserHTTP `yaml:"http"`
}

type UserHTTP struct {
	Timeout configlib.Duration `yaml:"timeout"`
}
