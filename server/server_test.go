package server

import (
	"github.com/faraway/wordofwisdom/config"
	"testing"
)

type (
	mockConfig struct{}
)

func newMockConfig() config.Config {
	return mockConfig{}
}

func (m mockConfig) Port() string {
	return ":8001"
}

func (m mockConfig) Protocol() string {
	return "tcp"
}

func TestServer_Run(t *testing.T) {
	srv := New(newMockConfig())

	err := srv.Run()
	if err != nil {
		t.Fatalf("server should start without errors got err: %v", err)
	}
}
