package protocol

import "encoding/json"

type (
	ClientMessage struct {
		Controller string          `json:"controller"`
		Message    json.RawMessage `json:"message"`
	}

	Action struct {
		Action  string      `json:"action"`
		Message interface{} `json:"message"`
	}

	Error struct {
		Action string `json:"action"`
		Error  string `json:"error"`
	}

	// Client -> Server
	ProveController struct {
		Suggest string `json:"suggest"`
	}

	// Server -> Client
	ChallengeAction struct {
		Challenge string `json:"challenge"`
	}

	ProveAction struct {
		Success bool
		Message string
	}
)
