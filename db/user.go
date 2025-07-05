package db

import (
	"context"
	"errors"
	"pgx-sqlc/db/sqlc"
)

// convenience function for use before UserDatabase can be created.
func GetUser(ctx context.Context, pgdb *Database, id string) (sqlc.User, error) {
	uuid, err := ReadUUID(id)
	if err != nil {
		panic(err)
	}
	return pgdb.Sqlc.GetUser(ctx, uuid)
}

// convenience function for use before UserDatabase can be created.
func NewUser(ctx context.Context, pgdb *Database, name, email string) (sqlc.User, error) {
	id, err := MakeUUIDv7()
	if err != nil {
		return sqlc.User{}, err
	}
	return pgdb.Sqlc.CreateUser(ctx, sqlc.CreateUserParams{
		ID:    id,
		Name:  name,
		Email: email,
	})
}

// Another wrapper.
// TODO: Not convinced this is the right abstraction. But surely I'll be doing most things on the
// behalf of a single user.
//
// NOTE: Prefer to pass this as an arg from within handlers, rather than add methods. I don't like
// implicit dependencies.
type UserDatabase struct {
	User *sqlc.User
	DB   *Database
}

func NewUserDatabase(user *sqlc.User, pgdb *Database) (UserDatabase, error) {
	if user == nil {
		return UserDatabase{}, errors.New("missing sqlc.User!")
	}
	if pgdb == nil {
		return UserDatabase{}, errors.New("missing Database!")
	}
	return UserDatabase{
		User: user,
		DB:   pgdb,
	}, nil
}
