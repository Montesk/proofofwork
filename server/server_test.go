package server

import (
	"github.com/Montesk/proofofwork/config"
	"github.com/Montesk/proofofwork/core/logger"
	"github.com/Montesk/proofofwork/router"
	"github.com/Montesk/proofofwork/sessioner"
	"testing"
	"time"
)

const (
	defaultTestWaitTime = 50 * time.Millisecond
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
	return 10 * time.Minute
}

func (m mockConfig) POWClients() int {
	return 0
}

func (m mockConfig) LogLevel() logger.Level {
	return logger.ErrorLevel
}

func newMockServer(port string) Server {
	ch := make(chan struct{})
	return New(newMockConfig(port), router.NewNull(), sessioner.NewMemory(ch), ch, logger.NewNull())
}

func TestServer_Run(t *testing.T) {
	srv := newMockServer(":8010")

	err := srv.Run()
	if err != nil {
		t.Fatalf("server should start without errors got err: %v", err)
	}

	srv.Close()
}

func TestServer_Listen(t *testing.T) {
	srv := newMockServer(":8011")

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
	srv := newMockServer(":8012")

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
