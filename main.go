package main

import (
	"fmt"
	"github.com/Montesk/proofofwork/cmd"
	"github.com/Montesk/proofofwork/config"
	"github.com/Montesk/proofofwork/core/logger"
	"github.com/Montesk/proofofwork/router"
	"github.com/Montesk/proofofwork/server"
	"github.com/Montesk/proofofwork/sessioner"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	args := cmd.NewFlagCmd()

	log := logger.NewBase(args.LogLevel())

	err := Run(config.NewBase(args.Protocol(), fmt.Sprintf(":%d", args.Port()), args.ReadTimeout(), args.POWClients(), args.LogLevel()), log)
	if err != nil {
		log.Fatalf("server run error err: %v", err)
	}

	log.Info("server shutdown")
}

func Run(cfg config.Config, log logger.Logger) error {
	srv := server.New(cfg, router.NewControllers(log), sessioner.NewMemory(), log)

	err := srv.Run()
	if err != nil {
		return err
	}

	log.Infof("tcp server start on port %s", cfg.Port())

	defer func() {
		err := srv.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Infof("listening...")

	return srv.Listen()
}
