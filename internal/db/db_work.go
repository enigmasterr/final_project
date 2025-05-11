package main

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ctx := context.TODO()

	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.PingContext(ctx)
	if err != nil {
		panic(err)
	}
}
