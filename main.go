package main

import (
	"fmt"


	"github.com/brm/api"
	
	db "github.com/brm/db/sqlc"
	"github.com/brm/logger"
	"github.com/brm/utils"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)


func main() {

	fmt.Println("Starting server...")

	config := utils.NewConfig()

	config.LoadEnvConfig()

	config.SanityCheck()

	con, err := db.NewPostgresDBConnection(*config)
	if err != nil {
		logger.Fatal("cannot open db connection:", zap.Error(err))
	}

	

	runDBMigration("file://db/migration", config.DBSource)

	store := db.NewStore(con)
	
	server, err := api.NewServer(config, store)
	if err != nil {
		logger.Fatal("cannot start server:", zap.Error(err))
	}

	err = server.Start()
	if err != nil {
		logger.Fatal("cannot start server:", zap.Error(err))
	}
}


func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		logger.Fatal("cannot create new migrate instance", zap.Error(err))
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Fatal("failed to run migrate up", zap.Error(err))
	}

	logger.Info("db migrated successfully")
}
