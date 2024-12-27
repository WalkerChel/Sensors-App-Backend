package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sensors-app/server"
	"syscall"
)

func Run() {
	httpServer := server.New(nil)
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
		log.Println("server was successfully stopped")
	}
}
