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

type FromConfig struct {
	Email string `yaml:"email"`
	Name  string `yaml:"name"`
}

type Mailer struct {
	Key  string     `yaml:"key"`
	From FromConfig `yaml:"from"`
}

type Config struct {
	Server Server `yaml:"server"`
	Mailer Mailer `yaml:"mailer"`
}
