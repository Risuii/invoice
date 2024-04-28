package main

import (
	"context"
	"log"
	"os"

	"github.com/Risuii/invoice/migration"
	"github.com/Risuii/invoice/src/app"
)

func main() {
	ctx := context.Background()

	app.Init(ctx)

	args := os.Args
	if len(args) < 2 {
		log.Fatal("Missing args. args: [up | rollback]")
	}

	migrationSvc, err := migration.New(ctx, app.Config().Postgres)
	if err != nil {
		log.Fatal("Failed to initiate migration", err)
	}

	switch args[1] {
	case "up":
		migrationSvc.Up(ctx)
	case "rollback":
		migrationSvc.Rollback(ctx)
	default:
		log.Fatal("Invalid migration command")
	}
}
