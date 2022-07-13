package protocol

type (
	Client struct {
		Controller string `json:"controller"`
		Message    string `json:"message"`
	}
)
