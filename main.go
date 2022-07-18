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

	l := logger.NewBase(logger.InfoLevel)
	args := cmd.NewFlagCmd(l)

	err := Run(config.NewBase(args.Protocol(), fmt.Sprintf(":%d", args.Port()), args.ReadTimeout(), args.POWClients(), args.LogLevel()), l)
	if err != nil {
		l.Fatalf("server run error err: %v", err)
	}

	l.Info("server shutdown")
}

func Run(cfg config.Config, log logger.Logger) error {
	srv := server.New(cfg, router.New(log), sessioner.New(), log)

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
