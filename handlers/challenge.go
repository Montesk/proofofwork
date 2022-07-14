package handlers

import (
	"fmt"
	"github.com/faraway/wordofwisdom/session"
	"log"
)

const (
	ChallengeController = "challenge"
	ChallengeAction     = "challenge"
)

func ChallengeHandler(ses session.Session) {
	err := ses.Send(struct {
		Action    string `json:"action"`
		Challenge string `json:"challenge"`
	}{
		Action:    ChallengeAction,
		Challenge: fmt.Sprintf("Take your challenge %v !", ses.ClientId()),
	})

	if err != nil {
		log.Fatal(err)
	}
}
