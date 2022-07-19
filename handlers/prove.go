package handlers

import (
	"github.com/Montesk/proofofwork/protocol"
	"github.com/Montesk/proofofwork/session"
)

func (h *handlers) ProveHandler(ses session.Session, msg protocol.ProveControllerMsg) {
	success := h.pow.Prove(ses.ClientId(), msg.Suggest)

	h.log.Debugf("received prove attempt from client %s msg %s", ses.ClientId(), msg.Suggest)

	var err error
	if success {
		err = ses.Send(protocol.ProveAction, protocol.ProveActionMsg{
			Success: true,
			Message: h.book.RandomQuote(),
		})
	} else {
		err = ses.Send(protocol.ProveAction, protocol.ProveActionMsg{
			Success: false,
			Message: "try again",
		})
	}

	if err != nil {
		h.log.Errorf("error sending message to client %s err %v", ses.ClientId(), err)
	}
}
