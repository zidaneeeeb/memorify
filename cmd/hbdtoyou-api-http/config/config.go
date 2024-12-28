package config

import "time"

type Config struct {
	Server     Server                `yaml:"server"`
	PostgreSQL map[string]PostgreSQL `yaml:"postgresql"`
	Encryption Encryption            `yaml:"encryption"`
	User       User                  `yaml:"user"`
	Content    Content               `yaml:"content"`
	Template   Template              `yaml:"template"`
	Payment    Payment               `yaml:"payment"`
}

type Server struct {
	Port        int `yaml:"port"`
	MetricsPort int `yaml:"metrics_port"`
}

type PostgreSQL struct {
	ConnectionString  string        `yaml:"connection_string"`
	ConnectionTimeout time.Duration `yaml:"connection_timeout"`
}

// Followings are the known PostgreSQL key on config
const (
	PostgreSQLTenant string = "tenant"
)

type Encryption struct {
	Key string `yaml:"key"`
}
