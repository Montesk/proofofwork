package handlers

import (
	"github.com/Montesk/wordofwisdom/protocol"
	"github.com/Montesk/wordofwisdom/session"
	"log"
)

const (
	ChallengeController = "challenge"
	ChallengeAction     = "challenge"
)

func (h *handlers) ChallengeHandler(ses session.Session, _ any) {
	challenge, err := h.pow.Generate(ses.ClientId())
	if err != nil {
		log.Print("error generating message to client ", err)
	}

	err = ses.Send(ChallengeAction, protocol.ChallengeAction{
		Challenge: challenge,
	})
	if err != nil {
		log.Print("error sending message to client ", err)
	}

}
