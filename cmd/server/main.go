package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"shortener/migrations"
	"syscall"
	"time"

	"shortener/internal/config"
	"shortener/internal/logger"
	"shortener/internal/repository"
	"shortener/internal/service"

	"github.com/wb-go/wbf/dbpg/pgx-driver"
	"github.com/wb-go/wbf/ginext"
	wbflogger "github.com/wb-go/wbf/logger"
	"github.com/wb-go/wbf/zlog"
)

func main() {
	zlog.InitConsole()

	cfg := config.Load()
	log := wbflogger.NewZerologAdapter("shoretner", "config.env")

	pg, err := pgxdriver.New(cfg.PostgresDSN, log,
		pgxdriver.MaxPoolSize(int32(cfg.DBMaxConns)),
		pgxdriver.MaxConnAttempts(cfg.DBConnAttempts),
		pgxdriver.BaseRetryDelay(cfg.DBRetryDelay),
		pgxdriver.MaxRetryDelay(cfg.DBMaxRetryDelay),
	)
	if err != nil {
		logger.Fatal("failed to connect to database", "error", err)
	}
	defer pg.Close()

	sqlDB, err := sql.Open("postgres", cfg.PostgresDSN)
	if err != nil {
		logger.Fatal("failed to open sql db for migrations", "error", err)
	}
	if err := migrations.Apply(sqlDB); err != nil {
		logger.Fatal("failed to apply migrations", "error", err)
	}
	sqlDB.Close()

	repo := repository.NewPostgresRepository(pg)

	svc := service.New(repo, cfg.BaseURL)

	router := ginext.New(cfg.GinMode)
	router.Use(ginext.Logger())
	router.Use(ginext.Recovery())

	router.GET("/", func(c *ginext.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	api := router.Group("/api")
	{
		api.POST("/shorten", svc.ShortenHandler)
		api.GET("/analytics/:short_url", svc.AnalyticsHandler)
	}

	router.GET("/s/:short_url", svc.RedirectHandler)

	router.LoadHTMLGlob("web/*.html")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		logger.Info("starting server", "port")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed", "error", err)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown error", "error", err)
	}
}
