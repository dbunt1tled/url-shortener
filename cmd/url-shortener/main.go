package main

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go_first/internal/config"
	"go_first/internal/lib/http-server/router"
	"go_first/internal/lib/logger"
	"go_first/storage/mysql"
	"net/http"
	"os"
)

func main() {
	var locale *i18n.Localizer
	cfg := config.MustLoadConfig()
	log := config.SetupLogger(cfg.Env, cfg.Debug)
	bundle := config.SetupLocale()
	//fmt.Println(locale.MustLocalize(&i18n.LocalizeConfig{
	//	MessageID: "already_exists",
	//	TemplateData: map[string]interface{}{
	//		"first_name":  "John",
	//		"second_name": "Doe",
	//	},
	//}))
	//fmt.Println(loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "already_exists"}))
	storage, err := mysql.Connection(cfg.DatabaseDSN)
	defer mysql.ConnectionClose(storage)
	if err != nil {
		log.Error("Error storage", logger.Error(err))
		os.Exit(1)
	}
	log.Debug("Start listening on address: " + cfg.HTTPServer.Address)
	server := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router.NewRouter(storage, cfg, locale, bundle, log),
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Error("Error listening on address: " + cfg.HTTPServer.Address)
		os.Exit(1)
	}
}
