package main

import (
	"context"
	"fmt"
	"os"

	"pgx-sqlc/db/sqlc"
)

func main() {
	ctx := context.Background()
	db_url := "user=bth database=testdb"

	// NOTE: this is the single connection method:
	// conn, err := pgx.Connect(ctx, db_url)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	// 	os.Exit(1)
	// }
	// defer conn.Close(ctx)
	// q := sqlc.New(conn)

	db := sqlc.NewDatabase(ctx, db_url)

	newacc, err := db.Query.CreateAccount(ctx, sqlc.CreateAccountParams{
		First: "Bob",
		Last:  "Jones",
		Email: "bob@j.com",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(newacc.ID)

	acc, err := db.Query.GetAccount(ctx, newacc.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "GetAuthor failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(acc.Email)

}
