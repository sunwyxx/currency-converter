package main

import (
	"context"
	_ "fmt"
	"github.com/sunwyx/currency-converter/internal/config"
	"github.com/sunwyx/currency-converter/internal/handler"
	"github.com/sunwyx/currency-converter/internal/infrastructure/api_client"
	"github.com/sunwyx/currency-converter/internal/infrastructure/redis_client"
	"github.com/sunwyx/currency-converter/internal/service"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	cfg:=config.MustLoad()
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level:slog.LevelDebug}),
			)
	ctx, _:= context.WithTimeout(context.Background(), 5*time.Second)
	cache, err := redis_client.NewCache(ctx, cfg.Redis)
	if err != nil{
		log.Error("failed to init redis","err", err)
		os.Exit(1)
	}
	api:= api_client.NewClient(cfg.API)
	conv := service.NewConverter(cache,  api, log)
	h := handler.NewHandler(conv)
	mux := http.NewServeMux()
	mux.HandleFunc("/convert", h.Convert)

	srv := &http.Server{
		Addr: ":8080",
		Handler:mux,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Info("server is running", "addr", srv.Addr)

	if err:= srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("server failed", "err", err)
		os.Exit(1)
	}
}