package config

import (
	"time"

	"xorm.io/xorm/log"
)

type EnvVariable string

type Authentication struct {
	AccessTokenExpiration time.Duration `yaml:"accessTokenExpirationMinutes"`
	SecretKey             string        `yaml:"secretKey"`
}

type Server struct {
	Port           int      `yaml:"port"`
	AllowedOrigins []string `yaml:"allowedOrigins"`
	Timeout        Timeout  `yaml:"timeout"`
}

type Timeout struct {
	// seconds
	Idle     time.Duration `yaml:"idle"`
	Read     time.Duration `yaml:"read"`
	Write    time.Duration `yaml:"write"`
	Shutdown time.Duration `yaml:"shutdown"`
}

type Logger struct {
	Level log.LogLevel `yaml:"level"`
}

type Database struct {
	Name     string `yaml:"name"`
	Ip       string `yaml:"ip"`
	Port     int    `yaml:"port"`
	Schema   string `yaml:"schema"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Logger   Logger `yaml:"logger"`
}

type Notification struct {
	Address string `yaml:"address"`
}

type Frontend struct {
	Url           string `yaml:"url"`
	ResetPassword string `yaml:"resetPassword"`
}

type Config struct {
	Server         Server         `yaml:"server"`
	Database       Database       `yaml:"database"`
	Authentication Authentication `yaml:"authentication"`
	Notification   Notification   `yaml:"notification"`
	Frontend       Frontend       `yaml:"frontend"`
}
