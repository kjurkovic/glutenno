package config

import (
	"time"

	"xorm.io/xorm/log"
)

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

type Authorization struct {
	Address string `yaml:"address"`
}

type Config struct {
	Server        Server        `yaml:"server"`
	Database      Database      `yaml:"database"`
	Authorization Authorization `yaml:"authorization"`
}
