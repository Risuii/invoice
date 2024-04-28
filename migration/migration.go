package migration

import (
	"context"
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	pg "github.com/Risuii/frs-lib/postgres"
	"github.com/Risuii/invoice/src/app"
)

const (
	migrateLogIdentifier = "payduct"
)

type MigrationService interface {
	Up(context.Context) error
	Rollback(context.Context) error
	Version(context.Context) (int, bool, error)
}

type migrationService struct {
	driver  database.Driver
	migrate *migrate.Migrate
}

func New(ctx context.Context, cfg app.Postgres) (MigrationService, error) {
	pgCfg := pg.PostgresConfig{
		ConnectionUrl: cfg.ConnURI,
	}

	sqlxDB, err := pg.InitSQLX(ctx, pgCfg)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	databaseInstance, err := postgres.WithInstance(sqlxDB.DB, &postgres.Config{})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	migrate, err := migrate.NewWithDatabaseInstance("file://migration/sql",
		migrateLogIdentifier, databaseInstance)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return migrationService{
		driver:  databaseInstance,
		migrate: migrate,
	}, nil
}

func (s migrationService) Up(ctx context.Context) error {
	currVersion, _, err := s.Version(ctx)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Running migration from version: %d", currVersion)
	if err := s.migrate.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println(err)
			return nil
		}
		log.Println(err)
		return err
	}

	currVersion, _, _ = s.Version(ctx)
	log.Println("current version: ", currVersion)
	log.Println(err)
	return nil
}

func (s migrationService) Rollback(ctx context.Context) error {
	currVersion, _, err := s.Version(ctx)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println(currVersion)

	if err := s.migrate.Steps(-1); err != nil {
		log.Println(err)
		return err
	}

	currVersion, _, _ = s.Version(ctx)
	log.Println("current version: ", currVersion)

	return nil
}

func (s migrationService) Version(ctx context.Context) (int, bool, error) {
	currVersion, dirty, err := s.driver.Version()
	if err != nil {
		log.Println(err)
		return 0, false, err
	}
	return currVersion, dirty, nil
}
