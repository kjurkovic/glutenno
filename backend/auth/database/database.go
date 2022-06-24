package database

import (
	"auth/config"
	"auth/models"
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

		// TODO: add sql logger
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d TimeZone=Europe/Zagreb",
			config.Ip,
			config.Username,
			config.Password,
			config.Name,
			config.Port,
		)

		// attempt to connect to db 3 times - fail after that
		var connection *gorm.DB
		var err error
		retries := 3
		for retries > 0 {
			connection, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
				NamingStrategy: schema.NamingStrategy{
					TablePrefix: fmt.Sprintf("%s.", config.Schema),
				},
			})
			if err == nil {
				break
			}
			time.Sleep(5 * time.Second)
			retries--
		}

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

	db.Db.AutoMigrate(&models.User{})
}
