package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"teach-stack/cmd/api"
	"teach-stack/configs"
	"teach-stack/db"
	"teach-stack/utils"
)

func main() {
	utils.SetupLogger(configs.Envs.Env)

	slog.Debug("Check connection string to database", slog.String("connectionString", configs.Envs.DatabaseConnString))
	db, err := db.NewPostgresStorage(configs.Envs.DatabaseConnString)
	if err != nil {
		slog.Error("failed to create a new db instance", utils.Err(err))
		return
	}
	initStorage(db)

	server := api.NewApi(fmt.Sprintf(":%s", configs.Envs.ApiPort), db)
	err = server.Run()
	if err != nil {
		slog.Error("failed to run server", utils.Err(err))
		return
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		slog.Error("cannot verify connection with database", utils.Err(err))
	}

	slog.Info("Database successfully connected!")
}
