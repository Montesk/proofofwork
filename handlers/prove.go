package handlers

import (
	"github.com/Montesk/proofofwork/protocol"
	"github.com/Montesk/proofofwork/session"
	"log"
)

const (
	ProveController = "prove"
	ProveAction     = "prove"
)

func (h *handlers) ProveHandler(ses session.Session, msg protocol.ProveController) {
	log.Printf("recieved prove attempt from client %s msg %s", ses.ClientId(), msg.Suggest)

	success := h.pow.Prove(ses.ClientId(), msg.Suggest)

	var err error
	if !success {
		err = ses.Send(ProveAction, protocol.ProveAction{
			Success: false,
			Message: "try again",
		})
	} else {
		err = ses.Send(ProveAction, protocol.ProveAction{
			Success: true,
			Message: h.book.RandomQuote(),
		})
	}

	if err != nil {
		log.Print("error sending message to client ", err)
	}

}
