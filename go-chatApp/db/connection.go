package db

import (
	"context"
	"go-chatapp/config"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context) (*pgxpool.Pool, error) {
	cfg := config.EnvConfig()

	config, err := pgxpool.ParseConfig(cfg.Database_URL)
	if err != nil {
		log.Printf("Database url is wrong")
		return nil, err
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Hour

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	log.Println("connection build successfully")

	return pool, nil

}
