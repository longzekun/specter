package db

import (
	"time"

	"github.com/longzekun/specter/server/config"
	"github.com/longzekun/specter/server/db/models"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newDBClient() *gorm.DB {
	dbConfig := config.GetServerConfig().DatabaseConfig
	var dbClient *gorm.DB
	switch dbConfig.Dialect {
	case config.Sqlite:
		dbClient = sqliteClient(&dbConfig)
	case config.MySQL:
		dbClient = mySQLClient(&dbConfig)
	case config.Postgres:
		dbClient = postgresClient(&dbConfig)
	default:
		panic("unsupported database dialect")
	}

	if dbClient == nil {
		panic("failed to connect database")
	}

	var allDBModels []interface{} = append(make([]interface{}, 0),
		&models.Certificate{},
	)
	var err error
	for _, model := range allDBModels {
		err = dbClient.AutoMigrate(model)
		if err != nil {
			zap.S().Fatal(err)
		}
	}

	// Get generic database object sql.DB to use its functions
	sqlDB, err := dbClient.DB()
	if err != nil {
		zap.S().Fatal(err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	return dbClient
}

func postgresClient(dbConfig *config.DatabaseConfig) *gorm.DB {
	dsn, err := dbConfig.DSN()
	if err != nil {
		panic(err)
	}
	dbClient, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		Logger:      getGormLogger(dbConfig),
	})
	if err != nil {
		panic(err)
	}
	return dbClient
}

func mySQLClient(dbConfig *config.DatabaseConfig) *gorm.DB {
	dsn, err := dbConfig.DSN()
	if err != nil {
		panic(err)
	}
	dbClient, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		Logger:      getGormLogger(dbConfig),
	})
	if err != nil {
		panic(err)
	}
	return dbClient
}
