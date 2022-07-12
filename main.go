package main

import (
	"github.com/faraway/wordofwisdom/server"
	"log"
)

func main() {
	srv := server.New("tcp", ":8001")

	err := RunApp(srv)
	if err != nil {
		log.Fatalf("server run error err: %v", err)
	}

	log.Printf("shutdown")
}

func RunApp(srv server.Server) error {
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
