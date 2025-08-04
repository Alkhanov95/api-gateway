package main

import (
	"context"
	"log"
	"os"

	"github.com/alkhanov95/api-gateway/handler"
	"github.com/alkhanov95/api-gateway/internal/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env не найден или не загружен")
	}

	dburl := os.Getenv("DATABASE_URL")
	if dburl == "" {
		log.Fatal("DATABASE_URL пуст")
	}

	conn, err := storage.GetConnect(context.Background(), dburl)
	if err != nil {
		log.Fatalf("ошибка подключения к БД: %v", err)
	}
	defer conn.Close()

	log.Println("✅ Подключение к базе прошло успешно")

	app := fiber.New()

	handler.SetupUserRoutes(app, conn)

	log.Fatal(app.Listen(":3000"))
}
