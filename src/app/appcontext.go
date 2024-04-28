package app

import (
	"context"
	"path/filepath"
	"runtime"

	frsI18n "github.com/Risuii/frs-lib/i18n"
	frsPostgres "github.com/Risuii/frs-lib/postgres"
	frsRedis "github.com/Risuii/frs-lib/redis"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
)

type appContext struct {
	db               *sqlx.DB
	requestValidator *validator.Validate
	redis            frsRedis.Redis
	cfg              *Configuration
}

var appCtx appContext
var appTransFile = func() string {
	_, f, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(filepath.Dir(f))

	// Return the project root directory path.
	return filepath.Join(basepath, "translation")
}()

func Init(ctx context.Context) error {

	cfg, err := InitConfig(ctx)
	if err != nil {
		return err
	}

	if err := frsI18n.Init(ctx, cfg.Translation.FilePath, appTransFile, cfg.Translation.DefaultLanguage); err != nil {
		panic(err)
	}

	db, err := frsPostgres.InitSQLX(ctx, frsPostgres.PostgresConfig{
		ConnectionUrl:      cfg.Postgres.ConnURI,
		MaxPoolSize:        cfg.Postgres.MaxPoolSize,
		MaxIdleConnections: cfg.Postgres.MaxIdleConnections,
		ConnMaxIdleTime:    cfg.Postgres.MaxIdleTime,
		ConnMaxLifeTime:    cfg.Postgres.MaxLifeTime,
	})
	if err != nil {
		return err
	}

	redis, err := frsRedis.InitRedis(ctx, cfg.Redis.Host, cfg.Redis.Password)
	if err != nil {
		return err
	}

	appCtx = appContext{
		db:               db,
		redis:            redis,
		requestValidator: validator.New(),
		cfg:              cfg,
	}

	return nil
}

func RequestValidator() *validator.Validate {
	return appCtx.requestValidator
}

func DB() *sqlx.DB {
	return appCtx.db
}

func Cache() frsRedis.Redis {
	return appCtx.redis
}

func Config() Configuration {
	return *appCtx.cfg
}
