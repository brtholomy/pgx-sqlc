package db

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func GetUUIDv7() (*pgtype.UUID, error) {
	googleuuid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	var id pgtype.UUID
	err = id.Scan(googleuuid.String())
	if err != nil {
		return nil, err
	}
	return &id, nil
}
