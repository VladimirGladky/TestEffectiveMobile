package main

import (
	"TestEffectiveMobile/internal/app"
	"TestEffectiveMobile/internal/config"
	"TestEffectiveMobile/pkg/logger"
	"context"
)

// @title Сервис подписок API
// @version 1.0.0
// @description REST-сервис для агрегации данных об онлайн-подписках пользователей.
// @host localhost:4047
// @BasePath /api/v1

func main() {
	ctx := context.Background()
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	ctx, err = logger.New(ctx)
	if err != nil {
		panic(err)
	}
	newApp := app.New(cfg, ctx)
	newApp.MustRun()
}
