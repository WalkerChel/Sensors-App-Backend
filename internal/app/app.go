package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sensors-app/configs"
	"sensors-app/db"

	"sensors-app/internal/api/handler"

	"sensors-app/internal/repository/repoPostgres"
	"sensors-app/internal/repository/repoRedis"
	"sensors-app/internal/service/serviceAuth"
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

	redisDB, err := db.NewRedisDB(context.Background(),
		cnf.Redis.Host,
		cnf.Redis.Port,
		cnf.Redis.Password,
		cnf.Redis.DB)

	if err != nil {
		log.Printf("redis db initialization err: %s", err)
		return
	}

	userRepo := repoPostgres.NewUserRepo(postgresDB)
	tokenRepo := repoRedis.NewTokenRepo(redisDB)

	authService := serviceAuth.NewAuthService(&tokenRepo, &userRepo)

	userHandlers := handler.NewUserHandlers(&authService)

	router := handler.Handlers{
		UserHandlers: userHandlers,
	}

	handlers := router.InitRoutes(cnf)

	httpServer := server.New(handlers)
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
