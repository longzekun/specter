package db

import (
	gosqlite "github.com/glebarez/sqlite"
	"github.com/longzekun/specter/server/config"
	"gorm.io/gorm"
)

func sqliteClient(dbConfig *config.DatabaseConfig) *gorm.DB {
	dsn, err := dbConfig.DSN()
	if err != nil {
		panic(err)
	}

	dbClient, err := gorm.Open(gosqlite.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		Logger:      getGormLogger(dbConfig),
	})
	if err != nil {
		panic(err)
	}
	return dbClient
}
