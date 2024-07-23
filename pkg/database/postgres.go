package database

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/adamnasrudin03/go-skeleton-chi/app/configs"
	"github.com/adamnasrudin03/go-skeleton-chi/app/models"
	"github.com/adamnasrudin03/go-skeleton-chi/pkg/seeders"
	"github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var (
	db  *gorm.DB
	err error
	mu  = &sync.Mutex{}
)

// SetupDbConnection is creating a new connection to our database
func SetupDbConnection(cfg *configs.Configs, logger *logrus.Logger) *gorm.DB {
	mu.Lock()
	defer mu.Unlock()
	logLevel := gormLogger.Silent
	if cfg.App.Env == "dev" {
		logLevel = gormLogger.Info
	}

	dbLogger := gormLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gormLogger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logLevel,    // Log level
			Colorful:      true,        // Disable color
		},
	)
	gormConfig := &gorm.Config{
		// enhance performance config
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger:                 dbLogger,
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DB.Host,
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.DbName,
		cfg.DB.Port)

	db, err = gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		logger.Panicf("Failed to create a connection to database , %v", err)
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Panicf("Failed to check connection to database , %v", err)
		return nil
	}
	if err := sqlDB.Ping(); err != nil {
		logger.Panicf("Failed to ping connection to database , %v", err)
		return nil
	}

	if cfg.DB.DbIsMigrate {
		//auto migration entity db
		db.AutoMigrate(
			&models.TeamMember{},
		)
	}

	go func(db *gorm.DB) {
		seeders.InitTeamMembers(db)
	}(db)

	logger.Info("Connection Database Success!")
	return db
}

// CloseDbConnection method is closing a connection between your app and your db
func CloseDbConnection(db *gorm.DB, logger *logrus.Logger) {
	dbSQL, err := db.DB()
	if err != nil {
		logger.Panicf("Failed to close connection from database, %v", err)
		return
	}

	dbSQL.Close()
}

func GetDB() *gorm.DB {
	return db
}
