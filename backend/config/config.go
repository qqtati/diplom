package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type HTTPConfig struct {
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
	Version string `yaml:"version"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
	SSLMode  string `yaml:"ssl_mode"`
}

type AuthConfig struct {
	Salt                 string `yaml:"salt"`
	AccessTokenLifetime  int    `yaml:"access_token_lifetime"`
	RefreshTokenLifetime int    `yaml:"refresh_token_lifetime"`
}

type RedisConfig struct {
	Host     string `yaml:"host" encrypted:"true"`
	Port     string `yaml:"port" encrypted:"true"`
	Password string `yaml:"password" encrypted:"true"`
}

type MailConfig struct {
	SMTPHost string `yaml:"smtp_host"`
	SMTPPort string `yaml:"smtp_port"`
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}

type Config struct {
	Postgres PostgresConfig `yaml:"postgres"`
	Server   HTTPConfig     `yaml:"server"`
	Auth     AuthConfig     `yaml:"auth"`
	Redis    RedisConfig    `yaml:"redis"`
	Mail     MailConfig     `yaml:"mail"`
}

func GenerateAppConfig() (*Config, error) {
	var cfg *Config
	var err error
	cfg, err = loadConfig("/app/config/config.yml")
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func loadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("file not found by %s", path)
	}
	yamlDecoder := yaml.NewDecoder(file)

	cfg := &Config{}
	err = yamlDecoder.Decode(cfg)
	if err != nil {
		return nil, fmt.Errorf("config: %v", err)
	}

	return cfg, nil
}
