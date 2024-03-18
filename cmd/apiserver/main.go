package main

import (
	"log/slog"
	"net/http"
	"os"

	_ "github.com/P1coFly/vk_movies/docs"
	"github.com/P1coFly/vk_movies/internal/config"
	"github.com/P1coFly/vk_movies/internal/http-server/handler"
	"github.com/P1coFly/vk_movies/internal/storage/postgresql"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Vk Movies API
// @version 1.0
// @description This is a RESTful API service for managing movies and actors
// @host localhost:8080
// @BasePath /api
func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting api-servies", "env", cfg.Env)
	log.Debug("cfg data", "data", cfg)

	storage, err := postgresql.New(cfg.Host_db)
	if err != nil {
		log.Error("failed to connect storage", "error", err)
		os.Exit(1)
	}
	log.Info("connect to db is successful", "host", cfg.Host_db)
	_ = storage

	http.HandleFunc("/api/actors", handler.ActorsHandler(storage))
	http.HandleFunc("/api/actor", handler.ActorHandler(storage, cfg))
	http.HandleFunc("/api/movie", handler.MovieHandler(storage, cfg))
	http.HandleFunc("/api/movies/byTitleFragment", handler.FindMoviesByTitleFragmentHandler(storage))
	http.HandleFunc("/api/movies/byActorNameFragment", handler.FindMoviesByActorNameFragmentHandler(storage))
	http.Handle("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

	log.Info("Сервер запущен на :8080")
	http.ListenAndServe(":8080", nil)

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "dev":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "prod":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log

}
