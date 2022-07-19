package main

import (
	"fmt"
	"github.com/Montesk/proofofwork/cmd"
	"github.com/Montesk/proofofwork/config"
	"github.com/Montesk/proofofwork/core/logger"
	"github.com/Montesk/proofofwork/pow/client"
	"github.com/Montesk/proofofwork/pow/networked"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

// network based suggestions by clients
func main() {
	rand.Seed(time.Now().UnixNano())

	args := cmd.NewFlagCmd()

	log := logger.NewBase(args.LogLevel())

	log.Printf("Running clients for proof of work...")

	done := make(chan struct{})

	wg := new(sync.WaitGroup)
	for i := 0; i < args.POWClients(); i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			log.Debugf("client N %d connected", idx)

			// server
			connectedClient := networked.New(config.NewBase(args.Protocol(), fmt.Sprintf("%d", args.Port()), args.ReadTimeout(), args.POWClients(), args.LogLevel()), log)

			challenge, err := connectedClient.Generate(strconv.Itoa(idx))
			if err != nil {
				log.Errorf("client %d failed to suggest challenge %v", idx, err)
				return
			}

			cl := client.New(strconv.Itoa(idx), connectedClient)

			tries, err := cl.Suggest(challenge, 0)
			if err != nil {
				log.Errorf("client %d failed to suggest challenge %v tries %d", idx, err, tries)
			} else {
				log.Infof("client %d suggested challenge in %d tries", idx, tries)
			}
		}(i)
	}
	wg.Wait()

	<-done
}
