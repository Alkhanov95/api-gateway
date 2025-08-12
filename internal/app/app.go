package app

import (
	"context"
	"os"

	"github.com/alkhanov95/api-gateway/internal/handler"
	"github.com/alkhanov95/api-gateway/internal/repository"
	"github.com/alkhanov95/api-gateway/internal/storage"
	"github.com/pkg/errors"

	"github.com/joho/godotenv"
	//slog
)

func Run() error {
	if err := godotenv.Load(); err != nil {
		return errors.Wrap(err, "env is not found")
	}

	dburl := os.Getenv("DATABASE_URL")
	if dburl == "" {
		return errors.New("DATABASE_URL empty")
	}

	conn, err := storage.GetConnect(context.Background(), dburl)
	if err != nil {
		return errors.Wrap(err, "connecting to DB")
	}
	defer conn.Close()

	repo := repository.NewUserRepo(conn)
	h := handler.New(repo)
	router := setupUserRoutes(h)

	if err := router.Listen(":3000"); err != nil {
		return errors.Wrap(err, "server listen") //message to docker/kuber that we shut with err
	}
	return nil
}
