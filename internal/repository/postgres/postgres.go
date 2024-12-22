package postgres

import (
	"context"
	"fmt"

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
	up := `CREATE TABLE IF NOT EXISTS users (
	tg_id int NOT NULL PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	
	firstname VARCHAR(255),
	lastname VARCHAR(255),
	username VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS log_request_response (
    id uuid PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	
	user_id int NOT NULL,
	request TEXT NOT NULL,
	response TEXT NOT NULL
);
`
	_, err := db.Exec(context.Background(), up)
	// files, err := os.ReadDir("./migrations")
	// if err != nil {
	// 	return err
	// }

	// for _, file := range files {
	// 	if !file.IsDir() && strings.HasSuffix(file.Name(), ".up.sql") {
	// 		content, err := os.ReadFile("./migrations/" + file.Name())
	// 		if err != nil {
	// 			return err
	// 		}
	// 		_, err = db.Exec(context.Background(), string(content))
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	return err
}
