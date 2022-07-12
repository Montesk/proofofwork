package main

import (
	"github.com/faraway/wordofwisdom/config"
	"github.com/faraway/wordofwisdom/server"
	"log"
)

func main() {
	err := RunApp(config.New("tcp", ":8001"))
	if err != nil {
		log.Fatalf("server run error err: %v", err)
	}

	log.Printf("shutdown")
}

func RunApp(cfg config.Config) error {
	srv := server.New(cfg)

	err := runServer(srv)
	if err != nil {
		return err
	}

	defer func() {
		err := srv.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Printf("tcp server start")

	return listenConnections(srv)
}

func runServer(srv server.Server) error {
	return srv.Run()
}

func listenConnections(srv server.Server) error {
	return srv.Listen()
}
