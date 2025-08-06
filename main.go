package main

import (
	"log/slog"
	"os"

	"github.com/alkhanov95/api-gateway/internal/app"
	//slog
)

func main() {
	if err := app.Run(); err != nil {
		slog.Error("app Run", slog.Any("error", err))
		os.Exit(1)
	}
}
