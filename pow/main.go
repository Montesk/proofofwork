package main

import (
	"github.com/Montesk/proofofwork/config"
	"github.com/Montesk/proofofwork/pow/client"
	"github.com/Montesk/proofofwork/pow/server"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const (
	clientsNum = 100
)

// network based suggestions by clients
func main() {
	rand.Seed(time.Now().UnixNano())

	log.Printf("Running clients for proof of work...")

	done := make(chan struct{})

	wg := new(sync.WaitGroup)
	for i := 0; i < clientsNum; i++ {
		go func(idx int) {
			wg.Add(1)
			defer wg.Done()

			log.Printf("client N %d connected", idx)

			// server
			connectedClient := server.New(config.New("tcp", "8001", 0))

			challenge, err := connectedClient.Generate(strconv.Itoa(idx))
			if err != nil {
				log.Fatal(err)
			}

			cl := client.New(strconv.Itoa(idx), connectedClient)

			tries, err := cl.Suggest(challenge, 0)
			if err != nil {
				log.Printf("client %d failed to suggest challenge %v tries %d", idx, err, tries)
			} else {
				log.Printf("client %d suggested challenge in %d tries", idx, tries)
			}

		}(i)
	}
	wg.Wait()

	<-done
}
