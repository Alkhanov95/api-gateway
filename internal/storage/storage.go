package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

func GetConnect(ctx context.Context, connstr string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(connstr)
	if err != nil {
		return nil, errors.Wrap(err, "не удалось распарсить строку подключения")
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "не удалось создать пул подключений")
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, errors.Wrap(err, "ping до базы не прошёл")
	}

	return pool, nil
}
