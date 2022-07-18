package main

import (
	"fmt"
	"github.com/Montesk/proofofwork/cmd"
	"github.com/Montesk/proofofwork/config"
	"github.com/Montesk/proofofwork/pow/client"
	"github.com/Montesk/proofofwork/pow/networked"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

// network based suggestions by clients
func main() {
	rand.Seed(time.Now().UnixNano())

	settings := cmd.NewFlagCmd()

	log.Printf("Running clients for proof of work...")

	done := make(chan struct{})

	wg := new(sync.WaitGroup)
	for i := 0; i < settings.POWClients(); i++ {
		go func(idx int) {
			wg.Add(1)
			defer wg.Done()

			log.Printf("client N %d connected", idx)

			// server
			connectedClient := networked.New(config.NewMockConfig(settings.Protocol(), fmt.Sprintf("%d", settings.Port()), settings.ReadTimeout(), settings.POWClients()))

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
