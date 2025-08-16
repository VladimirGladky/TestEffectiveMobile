package main

import (
	"TestEffectiveMobile/internal/app"
	"TestEffectiveMobile/internal/config"
	"TestEffectiveMobile/pkg/logger"
	"context"
)

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
