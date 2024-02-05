package main

import (
	"context"
	"embed"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"

	cfg "sdn_list/config"
	"sdn_list/infrastructure/db"
	"sdn_list/infrastructure/logger"
	"sdn_list/internal/clients"
	"sdn_list/internal/handlers"
	"sdn_list/internal/repositories"
	"sdn_list/internal/services"
	"sync/atomic"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	config := cfg.InitConfig()

	// logger
	logger, err := logger.InitLogger(cfg.ServiceName, config.Log.Level)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// db connection
	dbpool := db.NewPgxpool(ctx, config.DB, logger)
	defer dbpool.Close()

	uploadingInProcess := &atomic.Bool{}

	uploadAttemptsRepository := repositories.NewUploadAttemptsRepository(dbpool)
	sdnRepository := repositories.NewSdnRepository(dbpool)
	sdnXmlClient := clients.NewXmlClient("https://www.treasury.gov/ofac/downloads/sdn.xml")

	uploadService := services.NewUploadService(uploadAttemptsRepository, sdnRepository, sdnXmlClient)
	searchService := services.NewSearchService(sdnRepository)

	stateHandler := handlers.NewStateHandler(uploadService, uploadingInProcess)
	uploadHandler := handlers.NewUploadHandler(uploadService, uploadingInProcess)
	searchHandler := handlers.NewSearchHandler(searchService)

	http.HandleFunc("/state", stateHandler.Handle)
	http.HandleFunc("/update", uploadHandler.Handle)
	http.HandleFunc("/get_names", searchHandler.Handle)

	migrationsUp := func(w http.ResponseWriter, _ *http.Request) {
		goose.SetBaseFS(embedMigrations)
	
		if err := goose.SetDialect("postgres"); err != nil {
			io.WriteString(w, "Ошибка при выполнении миграций. Установка диалекта postgres. " + err.Error())
			return
		}
	
		db := stdlib.OpenDBFromPool(dbpool)
		if err := goose.Up(db, "migrations"); err != nil {
			io.WriteString(w, "Ошибка при выполнении миграций. " + err.Error())
			return
		}
		if err := db.Close(); err != nil {
			io.WriteString(w, "Ошибка при выполнении миграций. Закрытие соединения. " + err.Error())
			return
		}

		io.WriteString(w, "Миграции успешно применены")
	}

	http.HandleFunc("/migrations", migrationsUp)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
