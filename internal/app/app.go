package app

import (
	"context"

	"github.com/alkhanov95/api-gateway/config"
	"github.com/alkhanov95/api-gateway/internal/handler"
	"github.com/alkhanov95/api-gateway/internal/repository"
	"github.com/alkhanov95/api-gateway/internal/storage"
	"github.com/pkg/errors"
)

func Run() error {

	cfg, err := config.Parse()
	if err != nil {
		return errors.Wrap(err, "parsing config")
	}

	conn, err := storage.GetConnect(context.Background(), cfg.PGURL())
	if err != nil {
		return errors.Wrap(err, "connecting to DB")
	}
	defer conn.Close()

	repo := repository.NewUserRepo(conn)
	h := handler.New(repo)
	router := setupUserRoutes(h)

	if err := router.Listen(":" + cfg.App.Port); err != nil {
		return errors.Wrap(err, "server listen") //message to docker/kuber that we shut with err
	}
	return nil
}
