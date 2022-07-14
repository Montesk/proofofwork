package pow

import (
	"errors"
	"math/rand"
	"testing"
	"time"
)

func TestPow_GenerateAndProve(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	tests := map[string]struct {
		generatedClientId string
		requestClientId   string
		clientTries       int
		wantErr           bool
		err               error
	}{
		"generate and prove": {
			generatedClientId: "127.0.0.1:83887",
			requestClientId:   "127.0.0.1:83887",
		},
		"out of tries limit - err": {
			generatedClientId: "127.0.0.1:83887",
			requestClientId:   "127.0.0.1:83887",
			clientTries:       nonceMax + 100,
			wantErr:           true,
			err:               ErrTooManyTries,
		},
		"no block was registered for the user - err": {
			generatedClientId: "127.0.0.1:83887",
			requestClientId:   "83.88.10.100",
			wantErr:           true,
			err:               ErrTooManyTries,
		},
	}

	for title, test := range tests {
		t.Run(title, func(t *testing.T) {
			service := New()

			challenge, _ := service.Generate(test.generatedClientId)

			client := newMockClient(test.requestClientId, service)

			tries, err := client.Suggest(challenge, test.clientTries)
			if err != nil {
				if !test.wantErr {
					t.Errorf("expect no error got %v", err)
					return
				}

				if !errors.Is(err, test.err) {
					t.Errorf("expect err %v got err %v", test.wantErr, test.err)
				}
			}

			if test.wantErr {
				return
			}

			t.Logf("client suggested nonce in %d tries", tries)
		})
	}
}
