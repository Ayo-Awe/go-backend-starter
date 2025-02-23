package main

import (
	"cmp"
	"context"
	"flag"
	"log"
	"os"

	_ "github.com/ayo-awe/go-backend-starter/internal/db/migrations"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	flags := flag.NewFlagSet("goose", flag.ExitOnError)
	defaultDSN := cmp.Or(os.Getenv("GOOSE_DBSTRING"), "postgres://postgres:postgres@localhost:5432/go_starter?sslmode=disable")
	defaultDir := cmp.Or(os.Getenv("GOOSE_MIGRATION_DIR"), "./internal/db/migrations")

	dir := flags.String("dir", defaultDir, "directory with migration files")
	dsn := flags.String("dsn", defaultDSN, "databse dsn")

	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatalf("failed to parse flags: %v", err)
	}

	args := flags.Args()

	// Usage: goose -dir <directory> -dsn <database_url> COMMAND
	// command is the only required argument
	if len(args) < 1 {
		flags.Usage()
		return
	}

	command := args[0]
	cmdArgs := args[1:]

	db, err := goose.OpenDBWithDriver("pgx", *dsn)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v", err)
	}

	ctx := context.Background()
	if err := goose.RunContext(ctx, command, db, *dir, cmdArgs...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}
