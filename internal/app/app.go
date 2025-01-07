package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sensors-app/configs"
	"sensors-app/db"

	"sensors-app/internal/api/handler"
	"sensors-app/internal/service"

	"sensors-app/internal/repository/repoPostgres"
	"sensors-app/internal/repository/repoRedis"

	"sensors-app/server"
	"syscall"

	_ "github.com/lib/pq"
)

func Run() {
	cnf, err := configs.New()
	if err != nil {
		log.Printf("env config initialization err: %s", err)
		return
	}
	log.Print("Configuration variables was initialized successfully")

	postgresDB, err := db.NewPostgresDB(cnf.Postgres.Host,
		cnf.Postgres.Port,
		cnf.Postgres.User,
		cnf.Postgres.Name,
		cnf.Postgres.Password,
		cnf.Postgres.SSLMode)

	if err != nil {
		log.Printf("postgres db initialization err: %s", err)
		return
	}
	log.Print("Postgres DB was initialized successfully")

	redisDB, err := db.NewRedisDB(context.Background(),
		cnf.Redis.Host,
		cnf.Redis.Port,
		cnf.Redis.Password,
		cnf.Redis.DB)

	if err != nil {
		log.Printf("redis db initialization err: %s", err)
		return
	}
	log.Print("Redis DB was initialized successfully")

	userRepo := repoPostgres.NewUserRepo(postgresDB)
	log.Print("User repository was initialized successfully")

	tokenRepo := repoRedis.NewTokenRepo(redisDB)
	log.Print("Token repository was initialized successfully")

	regionRepo := repoPostgres.NewRegionsRepo(postgresDB)
	log.Print("Region repository was initialized successfully")

	authService := service.NewAuthService(&tokenRepo)
	log.Print("Authentication service was initialized successfully")

	userService := service.NewUserService(&userRepo)
	log.Print("User service was initialized successfully")

	regionService := service.NewRegionService(&regionRepo)
	log.Print("Region service was initialized successfully")

	userHandlers := handler.NewUserHandlers(&userService)
	log.Print("User handlers was initialized successfully")

	regionHandlers := handler.NewRegionService(&regionService)
	log.Print("Region handlers was initialized successfully")

	router := handler.Handlers{
		UserHandlers:   userHandlers,
		RegionHandlers: regionHandlers,
	}
	log.Print("Router struct was initialized successfully")

	handlers := router.InitRoutes(cnf, &authService)
	log.Print("Router was initialized successfully")

	httpServer := server.New(handlers,
		server.AddAddress(cnf.Server.Host, cnf.Server.Port))
	log.Print("Server was initialized successfully")
	httpServer.Start()

	appCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL, os.Interrupt)
	defer stop()

	select {
	case <-appCtx.Done():
		log.Print("server was manually stopped")
	case err := <-httpServer.Notify():
		log.Printf("server error: %s", err)
	}

	log.Print("closing server")
	if err := httpServer.Shutdown(); err != nil {
		log.Printf("error occured when closing server: %s", err)
	} else {
		log.Print("server was successfully stopped")
	}
}
