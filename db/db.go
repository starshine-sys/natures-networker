package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/starshine-sys/natures-networker/common"

	// pgx driver for migrations
	_ "github.com/jackc/pgx/v4/stdlib"
)

type DB struct {
	*pgxpool.Pool
}

func New(dsn string) (*DB, error) {
	err := runMigrations(dsn)
	if err != nil {
		return nil, err
	}

	pgxconf, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse config: %w", err)
	}
	pgxconf.ConnConfig.LogLevel = pgx.LogLevelWarn
	pgxconf.ConnConfig.Logger = zapadapter.NewLogger(common.Log.Desugar())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := pgxpool.ConnectConfig(ctx, pgxconf)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %w", err)
	}

	db := &DB{Pool: pool}

	return db, nil
}

//go:embed migrations
var fs embed.FS

func runMigrations(url string) (err error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: fs,
		Root:       "migrations",
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}

	if n != 0 {
		common.Log.Infof("Performed %v migrations!", n)
	}

	err = db.Close()
	return err
}
