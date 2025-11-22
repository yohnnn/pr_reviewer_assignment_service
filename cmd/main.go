package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/yohnnn/pr_reviewer_assignment_service/internal/app/handlers"
	"github.com/yohnnn/pr_reviewer_assignment_service/internal/app/router"
	"github.com/yohnnn/pr_reviewer_assignment_service/internal/config"
	"github.com/yohnnn/pr_reviewer_assignment_service/internal/logger"
	postgresRepo "github.com/yohnnn/pr_reviewer_assignment_service/internal/repository/postgres"
	"github.com/yohnnn/pr_reviewer_assignment_service/internal/services"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	log := logger.New("prod", cfg.Server.LogLevel)

	ctx := context.Background()
	poolConfig, err := pgxpool.ParseConfig(cfg.Postgres.ConnectionString())
	if err != nil {
		log.Error("failed to parse db config", "error", err)
		os.Exit(1)
	}

	poolConfig.MaxConns = 20
	poolConfig.MinConns = 5

	dbPool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Error("failed to create db pool", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		log.Error("failed to ping db", "error", err)
		os.Exit(1)
	}
	log.Info("connected to database via pgxpool")

	userRepo := postgresRepo.NewUserRepository(dbPool)
	teamRepo := postgresRepo.NewTeamRepository(dbPool)
	prRepo := postgresRepo.NewPullRequestRepository(dbPool)

	userService := services.NewUserService(userRepo, log)
	teamService := services.NewTeamService(teamRepo, log)
	prService := services.NewPullRequestService(prRepo, userRepo, log)

	userHandler := handlers.NewUserHandler(userService, prService, log)
	teamHandler := handlers.NewTeamHandler(teamService, log)
	prHandler := handlers.NewPullRequestHandler(prService, log)

	r := router.NewRouter(userHandler, teamHandler, prHandler)
	ginEngine := r.InitRoutes()

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: ginEngine,
	}

	go func() {
		log.Info("server started", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("forced shutdown", "error", err)
	}
	log.Info("server exited")
}
