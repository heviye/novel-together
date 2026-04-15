package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	App      AppConfig      `yaml:"app"`
	JWT      JWTConfig      `yaml:"jwt"`
	CORS     CORSConfig     `yaml:"cors"`
}

type DatabaseConfig struct {
	Driver   string `yaml:"driver"` // sqlite or postgres
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Path     string `yaml:"path"` // only for sqlite
}

type AppConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
}

type CORSConfig struct {
	Origins string `yaml:"origins"` // comma-separated list of allowed origins, "*" for all
}

var AppConf *Config

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	AppConf = &cfg
	return &cfg, nil
}

func (d *DatabaseConfig) DSN() string {
	switch d.Driver {
	case "sqlite":
		return d.Path
	case "postgres":
		return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
			d.Host, d.User, d.Password, d.Name, d.Port)
	default:
		return d.Path // default to sqlite
	}
}

func (d *DatabaseConfig) IsSQLite() bool {
	return d.Driver == "sqlite" || d.Driver == ""
}

func (d *DatabaseConfig) IsPostgres() bool {
	return d.Driver == "postgres"
}

func (c *CORSConfig) AllowedOrigin(origin string) bool {
	if c.Origins == "" || c.Origins == "*" {
		return true
	}
	for _, o := range strings.Split(c.Origins, ",") {
		if strings.TrimSpace(o) == origin {
			return true
		}
	}
	return false
}
