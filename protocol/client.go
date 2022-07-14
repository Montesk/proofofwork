package protocol

type (
	ClientMessage struct {
		Controller string `json:"controller"`
		Message    string `json:"message"`
	}

	Action struct {
		Action  string      `json:"action"`
		Message interface{} `json:"message"`
	}

	ChallengeAction struct {
		Challenge string `json:"challenge"`
	}
)
