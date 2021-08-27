package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yosa12978/WikiMD/internal/config"
	server "github.com/yosa12978/WikiMD/internal/web"
)

func Run() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	server.Run(cfg.Server.Port)

	out := make(chan os.Signal, 1)
	signal.Notify(out, syscall.SIGINT, syscall.SIGTERM)
	sig := <-out
	log.Printf("Programm stopped at %d, signal %s\n", time.Now().Unix(), sig.String())
}
