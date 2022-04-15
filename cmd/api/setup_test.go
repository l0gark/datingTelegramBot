package main

import (
	"context"
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
	app, _, err := initApp()
	if err != nil {
		log.Fatalf("could not create app %s", err.Error())
	}
	db = app.users.DB.(*pgxpool.Pool)

	_ = refreshLikesTable()
	_ = refreshUsersTable()

	return app
}

func refreshUsersTable() error {
	ctx := context.Background()
	conn, err := db.Acquire(ctx)
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Exec(ctx, "DELETE FROM users;")

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

	_, err = conn.Exec(ctx, "DELETE FROM likes;")

	if err != nil {
		log.Fatalf("Error truncating likes table: %s", err)
	}
	return nil
}
