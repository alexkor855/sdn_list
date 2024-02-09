package main

import (
	"context"
	"embed"
	"fmt"
	"io"
	"net"
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
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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

	mux := http.NewServeMux()
	mux.HandleFunc("/state", stateHandler.Handle)
	mux.HandleFunc("/update", uploadHandler.Handle)
	mux.HandleFunc("/get_names", searchHandler.Handle)

	migrationsUp := func(w http.ResponseWriter, _ *http.Request) {
		goose.SetBaseFS(embedMigrations)

		if err := goose.SetDialect("postgres"); err != nil {
			io.WriteString(w, "Ошибка при выполнении миграций. Установка диалекта postgres. "+err.Error())
			return
		}

		db := stdlib.OpenDBFromPool(dbpool)
		if err := goose.Up(db, "migrations"); err != nil {
			io.WriteString(w, "Ошибка при выполнении миграций. "+err.Error())
			return
		}
		if err := db.Close(); err != nil {
			io.WriteString(w, "Ошибка при выполнении миграций. Закрытие соединения. "+err.Error())
			return
		}

		io.WriteString(w, "Миграции успешно применены")
	}

	mux.HandleFunc("/migrations", migrationsUp)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Http.Port),
		Handler: mux,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return httpServer.ListenAndServe()
	})

	g.Go(func() error {
		<-ctx.Done()
		return httpServer.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		logger.Fatal("app exit reason: %s \n", zap.Error(err))
	}
}
