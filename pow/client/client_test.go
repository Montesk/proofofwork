package client

import (
	"errors"
	"github.com/Montesk/proofofwork/pow/pow"
	"github.com/Montesk/proofofwork/pow/service"
	"math/rand"
	"testing"
	"time"
)

func TestPow_GenerateAndProve(t *testing.T) {
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
			clientTries:       pow.NonceMax + 100,
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
			rand.Seed(time.Now().UnixNano())

			sv := service.New()

			challenge, err := sv.Generate(test.generatedClientId)
			if err != nil {
				t.Fatal(err)
			}

			cl := New(test.requestClientId, sv)

			tries, err := cl.Suggest(challenge, test.clientTries)
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
