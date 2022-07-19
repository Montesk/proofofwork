// Package provides description of DTO objects

package protocol

import "encoding/json"

// Actions
const (
	ChallengeAction = "challenge"
	ProveAction     = "prove"
	InfoAction      = "info"
)

// Controllers
const (
	ChallengeController = "challenge"
	ProveController     = "prove"
)

// DTO
type (
	ClientMessage struct {
		Controller string          `json:"controller"`
		Message    json.RawMessage `json:"message"`
	}

	Action struct {
		Action  string `json:"action"`
		Message any    `json:"message"`
	}

	Error struct {
		Action string `json:"action"`
		Error  string `json:"error"`
	}

	// Client -> Server
	ProveControllerMsg struct {
		Suggest string `json:"suggest"`
	}

	// Server -> Client
	ChallengeActionMsg struct {
		Challenge string `json:"challenge"`
	}

	ProveActionMsg struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	ServerInfoActionMsg struct {
		Info string `json:"info"`
	}
)
