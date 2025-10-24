package main

import (
	"fmt"
	"lumen/internal/config"
	"lumen/internal/databse"
	"lumen/internal/handler"
	"lumen/internal/model"
	"lumen/internal/repository"
	"lumen/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		panic(err)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	db, err := databse.ConnectSQLite()
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&model.Genre{})

	repo := repository.NewGenreRepo(db)
	service := service.NewGenreService(repo)
	validate := validator.New()

	h := handler.NewGenreHandler(service, validate, logger)

	r := chi.NewRouter()
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(handler.JSONContentType)
	r.Route("/api/genres", h.RegisterRoutes)

	addr := fmt.Sprintf(":%s", cfg.Port)
	logger.Info("Server listening", zap.String("port", cfg.Port))
	http.ListenAndServe(addr, r)
}
