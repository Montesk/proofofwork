package server

import (
	"github.com/faraway/wordofwisdom/config"
	"github.com/faraway/wordofwisdom/router"
	"testing"
	"time"
)

const (
	defaultTestWaitTime = 100 * time.Millisecond
)

type (
	mockConfig struct {
		port string
	}
)

func newMockConfig(port string) config.Config {
	return mockConfig{
		port: port,
	}
}

func (m mockConfig) Port() string {
	return m.port
}

func (m mockConfig) Protocol() string {
	return "tcp"
}

func (m mockConfig) ReadTimeout() time.Duration {
	return DefaultReadTimeout
}

func TestServer_Run(t *testing.T) {
	srv := New(newMockConfig(":8001"), router.Null())

	err := srv.Run()
	if err != nil {
		t.Fatalf("server should start without errors got err: %v", err)
	}

	srv.Close()
}

func TestServer_Listen(t *testing.T) {
	srv := New(newMockConfig(":8002"), router.Null())

	err := srv.Run()
	if err != nil {
		t.Fatal(err)
	}

	wait := make(chan struct{})
	go func() {
		time.Sleep(defaultTestWaitTime)
		wait <- struct{}{}
	}()

	go func() {
		defer srv.Close()

		err := srv.Listen()
		if err != nil {
			t.Errorf("server should start listening without error got err: %v", err)
			return
		}
	}()

	<-wait
}

func TestServer_Close(t *testing.T) {
	srv := New(newMockConfig(":8003"), router.Null())

	err := srv.Run()
	if err != nil {
		t.Fatal(err)
	}

	wait := make(chan struct{})
	go func() {
		time.Sleep(defaultTestWaitTime)
		wait <- struct{}{}
	}()

	go func() {
		err := srv.Close()
		if err != nil {
			t.Errorf("server should be closed without error got err: %v", err)
			return
		}
	}()

	<-wait
}
