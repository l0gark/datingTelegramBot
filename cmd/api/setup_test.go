package main

import (
	"context"
	"github.com/Eretic431/datingTelegramBot/internal/data/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
	"testing"
)

var (
	db *pgxpool.Pool
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func newTestApp() *application {
	mainConfig, err := getConfig()
	if err != nil {
		log.Printf("%e", err)
		return nil
	}
	sugaredLogger, _, err := newLogger(mainConfig)
	if err != nil {
		log.Printf("%e", err)
		return nil
	}
	postgresConfig := newPostgresConfig(mainConfig, sugaredLogger)
	_, _, err = postgres.NewPsqlPool(postgresConfig)
	if err != nil {
		log.Printf("%e", err)
		return nil
	}
	app, _, err := initApp()
	if err != nil {
		log.Fatalf("could not create app %s", err.Error())
	}
	_ = refreshUsersTable()
	_ = refreshLikesTable()

	return app
}

func refreshUsersTable() error {
	ctx := context.Background()
	conn, err := db.Acquire(ctx)
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Exec(ctx, "TRUNCATE TABLE users;")

	if err != nil {
		log.Fatalf("Error truncating users table: %s", err)
	}
	return nil
}

func refreshLikesTable() error {
	ctx := context.Background()
	conn, err := db.Acquire(ctx)
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Exec(ctx, "TRUNCATE TABLE likes;")

	if err != nil {
		log.Fatalf("Error truncating likes table: %s", err)
	}
	return nil
}
