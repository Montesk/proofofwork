package handlers

import (
	"github.com/Montesk/proofofwork/protocol"
	"github.com/Montesk/proofofwork/session"
)

func (h *handlers) ChallengeHandler(ses session.Session, _ any) {
	challenge, err := h.pow.Generate(ses.ClientId())
	if err != nil {
		h.log.Errorf("error generating message to client err %v", err)
		return
	}

	err = ses.Send(protocol.ChallengeAction, protocol.ChallengeActionMsg{
		Challenge: challenge,
	})
	if err != nil {
		h.log.Errorf("error sending message to client err %v", err)
	}

}
