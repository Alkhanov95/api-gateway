package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/alkhanov95/api-gateway/internal/handler"
	"github.com/alkhanov95/api-gateway/internal/storage"
	"github.com/pkg/errors"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	//slog
)

func Run() error {
	if err := godotenv.Load(); err != nil {
		return errors.Wrap(err, "env is not found")
	}

	dburl := os.Getenv("DATABASE_URL")
	if dburl == "" { //obyedenit?
		slog.Error("DATABASE_URL пуст")
		os.Exit(1)
	}

	conn, err := storage.GetConnect(context.Background(), dburl)
	if err != nil {
		slog.Error("ошибка подключения к БД: %v", slog.Any("error", err))
		os.Exit(1)
	}
	defer conn.Close() //wont work because of os.exit

	slog.Info("✅ Подключение к базе прошло успешно") //SLOG INFO ?

	app := fiber.New()

	h := &handler.Handle{DB: conn}
	SetupUserRoutes(app, h)

	if err := app.Listen(":3000"); err != nil {
		slog.Error("app run", slog.Any("error", err))
		os.Exit(1) //message to docker/kuber that we shut with err
	}
	return nil
}
