package handlers

import (
	"github.com/faraway/wordofwisdom/protocol"
	"github.com/faraway/wordofwisdom/session"
	"log"
)

const (
	ProveController = "prove"
	ProveAction     = "prove"
)

func (h *handlers) ProveHandler(ses session.Session, msg protocol.ProveController) {
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
