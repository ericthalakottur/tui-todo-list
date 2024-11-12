package main

import (
	"context"
	"database/sql"
	_ "embed"
	"log"

	"github.com/ericthalakottur/tui-todo-list/tasks"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var ddl string

type DBConnection struct {
	ctx     context.Context
	queries *tasks.Queries
}

func initializeDB() (*DBConnection, error) {
	dbConnection := DBConnection{
		ctx: context.Background(),
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return &DBConnection{}, err
	}
	log.Print("Connected to sqlite3 database")

	// create tables
	if _, err := db.ExecContext(dbConnection.ctx, ddl); err != nil {
		return &DBConnection{}, err
	}
	log.Print("Created table")

	dbConnection.queries = tasks.New(db)

	return &dbConnection, nil
}

func (dbConnection *DBConnection) CreateNewTask(taskName, completeBy string) error {
	log.Printf("Saving Task: %s\n", taskName)
	return dbConnection.queries.NewTask(
		dbConnection.ctx,
		tasks.NewTaskParams{
			Name:       taskName,
			CompleteBy: sql.NullString{String: completeBy, Valid: true},
		},
	)
}

func (dbConnection *DBConnection) GetIncompleteTasks() ([]tasks.Task, error) {
	return dbConnection.queries.GetIncompleteTasks(dbConnection.ctx)
}
