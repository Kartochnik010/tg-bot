package postgres

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(dsn string) (*pgxpool.Pool, error) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	err = autoMigrate(pool)
	if err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	return pool, nil
}

func autoMigrate(db *pgxpool.Pool) error {
	files, err := os.ReadDir("./migrations")
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".up.sql") {
			content, err := os.ReadFile("./migrations/" + file.Name())
			if err != nil {
				return err
			}
			_, err = db.Exec(context.Background(), string(content))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
