package main

import (
	"github.com/faraway/wordofwisdom/config"
	"github.com/faraway/wordofwisdom/router"
	"github.com/faraway/wordofwisdom/server"
	"github.com/faraway/wordofwisdom/sessioner"
	"log"
)

func main() {
	err := Run(config.New("tcp", ":8001", server.DefaultReadTimeout))
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
