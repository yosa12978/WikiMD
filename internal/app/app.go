package app

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yosa12978/WikiMD/internal/config"
	"github.com/yosa12978/WikiMD/internal/pkg/mongo"
	server "github.com/yosa12978/WikiMD/internal/web"
)

func Run() {
	rand.Seed(time.Now().Unix())
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	_, err = mongo.InitMongo(cfg)
	if err != nil {
		panic(err)
	}

	server.Run(cfg.Server.Port)

	out := make(chan os.Signal, 1)
	signal.Notify(out, syscall.SIGINT, syscall.SIGTERM)
	sig := <-out
	log.Printf("Programm stopped at %d, signal %s\n", time.Now().Unix(), sig.String())
}
