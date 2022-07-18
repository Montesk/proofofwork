package main

import (
	"fmt"
	"github.com/Montesk/proofofwork/cmd"
	"github.com/Montesk/proofofwork/config"
	"github.com/Montesk/proofofwork/router"
	"github.com/Montesk/proofofwork/server"
	"github.com/Montesk/proofofwork/sessioner"
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	settings := cmd.NewFlagCmd()

	err := Run(config.NewMockConfig(settings.Protocol(), fmt.Sprintf(":%d", settings.Port()), settings.ReadTimeout(), settings.POWClients()))
	if err != nil {
		log.Fatalf("server run error err: %v", err)
	}

	log.Printf("shutdown")
}

func Run(cfg config.Config) error {
	srv := server.New(cfg, router.New(), sessioner.New())

	err := srv.Run()
	if err != nil {
		return err
	}

	log.Printf("tcp server start")

	defer func() {
		err := srv.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Printf("listening...")

	return srv.Listen()
}
