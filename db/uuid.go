package db

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func ReadUUID(id string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	err := uuid.Scan(id)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return uuid, nil
}

func MakeUUIDv7() (pgtype.UUID, error) {
	googleuuid, err := uuid.NewV7()
	if err != nil {
		return pgtype.UUID{}, err
	}
	return ReadUUID(googleuuid.String())
}
