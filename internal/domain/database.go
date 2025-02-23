package domain

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Database interface {
	WithTransaction(context.Context, func(tx pgx.Tx) error) error
}
