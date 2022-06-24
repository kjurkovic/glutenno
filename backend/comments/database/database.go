package database

import (
	"comments/config"
	"comments/models"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Database struct {
	Db *gorm.DB
}

var once sync.Once
var (
	Instance *Database
)

func Get(config *config.Database, logger *log.Logger) *Database {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d TimeZone=Europe/Zagreb",
			config.Ip,
			config.Username,
			config.Password,
			config.Name,
			config.Port,
		)

		newLogger := gormLogger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			gormLogger.Config{
				SlowThreshold:             time.Second,     // Slow SQL threshold
				LogLevel:                  gormLogger.Info, // Log level
				IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,            // Disable color
			},
		)

		connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: newLogger,
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: fmt.Sprintf("%s.", config.Schema),
			},
		})

		if err != nil {
			logger.Fatal(err)
			return
		}
		Instance = &Database{Db: connection}
		migrations(Instance, config)
	})

	return Instance
}

func migrations(db *Database, config *config.Database) {
	if config.Schema != "public" {
		gormDb, _ := db.Db.DB()
		gormDb.Exec("CREATE SCHEMA IF NOT EXISTS " + config.Schema)
	}
	db.Db.AutoMigrate(&models.Comment{})
}
