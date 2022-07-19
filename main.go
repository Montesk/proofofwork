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
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	args := cmd.NewFlagCmd()

	log := logger.NewBase(args.LogLevel())

	shutdown, stop := make(chan struct{}), make(chan struct{})
	ss := sessioner.NewMemory(shutdown)

	go gracefulShutdown(ss, shutdown, stop, log)

	cfg := config.NewBase(args.Protocol(), fmt.Sprintf(":%d", args.Port()), args.ReadTimeout(), args.POWClients(), args.LogLevel())

	err := Run(cfg, ss, shutdown, log)
	if err != nil {
		log.Fatalf("server run error err: %v", err)
	}

	<-stop
}

func Run(cfg config.Config, ss sessioner.Sessioner, shutdown chan struct{}, log logger.Logger) error {
	srv := server.New(cfg, router.NewControllers(log), ss, shutdown, log)

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

// :WARNING: blocking function
func gracefulShutdown(ss sessioner.Sessioner, shutdown chan struct{}, stop chan struct{}, log logger.Logger) {
	defer func() {
		log.Infof("bye")
		stop <- struct{}{}
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	s := <-sig

	log.Infof("received signal %s", s.String())
	log.Infof("starting graceful shutdown")

	notifyServerShutdown(shutdown)
	go ss.Close()

	for {
		select {
		case <-time.After(10 * time.Second):
			log.Errorf("graceful shutdown timeout")
			return
		case <-shutdown:
			log.Infof("all clients disconnected")
			return
		}
	}
}

func notifyServerShutdown(shutdown chan struct{}) {
	shutdown <- struct{}{}
}
