package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go_first/internal/config"
	httpMiddlewares "go_first/internal/lib/http-server/middleware"
	"go_first/internal/lib/logger"
	"go_first/storage/mysql"
	"os"
	"time"
)

func main() {
	cfg := config.MustLoadConfig()
	log := config.SetupLogger(cfg.Env, cfg.Debug)
	storage, err := mysql.Connection(cfg.DatabaseDSN)
	defer mysql.ConnectionClose(storage)
	if err != nil {
		log.Error("Error storage", logger.Error(err))
		os.Exit(1)
	}
	//id, err := storage.CreateURL("https://www.google.com", "google")
	//if err != nil {
	//	log.Error("Error storage", logger.Error(err))
	//	os.Exit(1)
	//}

	//url, err := storage.GetURL(mysql.URLFilter{Alias: "google"})
	//fmt.Println(url)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Compress(5))
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(httpMiddlewares.Logger(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	log.Debug("Hello world Debug")
	fmt.Println(cfg)
}
