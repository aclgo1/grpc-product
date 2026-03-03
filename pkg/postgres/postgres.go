package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/aclgo/product/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect(ctx context.Context, c *config.Config) (*sqlx.DB, error) {
	conn, err := sqlx.Open(c.DatabaseDriver, c.DatabaseUrl)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Open: %v", err)
	}

	if err := conn.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("conn.PingContext: %v", err)
	}

	conn.SetMaxIdleConns(15)
	conn.SetMaxOpenConns(25)

	conn.SetConnMaxLifetime(time.Minute * 5)

	return conn, nil
}
