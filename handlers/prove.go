package handlers

import (
	"github.com/Montesk/proofofwork/protocol"
	"github.com/Montesk/proofofwork/session"
)

const (
	ProveController = "prove"
	ProveAction     = "prove"
)

func (h *handlers) ProveHandler(ses session.Session, msg protocol.ProveController) {
	success := h.pow.Prove(ses.ClientId(), msg.Suggest)

	h.log.Debugf("received prove attempt from client %s msg %s", ses.ClientId(), msg.Suggest)

	var err error
	if success {
		err = ses.Send(ProveAction, protocol.ProveAction{
			Success: true,
			Message: h.book.RandomQuote(),
		})
	} else {
		err = ses.Send(ProveAction, protocol.ProveAction{
			Success: false,
			Message: "try again",
		})
	}

	if err != nil {
		h.log.Errorf("error sending message to client %s err %v", ses.ClientId(), err)
	}
}
