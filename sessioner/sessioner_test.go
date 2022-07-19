package sessioner

import (
	"errors"
	"net"
	"testing"
	"time"
)

type (
	mockClient struct {
		RemoteAddr string
	}
	mockConn struct{}
)

func (m mockClient) Listen() error { return nil }
func (m mockClient) Close() error  { return nil }
func (m mockClient) Addr() string  { return m.RemoteAddr }

func (m mockConn) Read(b []byte) (n int, err error) {
	panic("implement me")
}

func (m mockConn) Write(b []byte) (n int, err error) {
	panic("implement me")
}

func (m mockConn) Close() error {
	panic("implement me")
}

func (m mockConn) LocalAddr() net.Addr {
	panic("implement me")
}

func (m mockConn) RemoteAddr() net.Addr {
	panic("implement me")
}

func (m mockConn) SetDeadline(t time.Time) error {
	panic("implement me")
}

func (m mockConn) SetReadDeadline(t time.Time) error {
	panic("implement me")
}

func (m mockConn) SetWriteDeadline(t time.Time) error {
	panic("implement me")
}

func newMockSessioner() Sessioner {
	return NewMemory(make(chan struct{}))
}

func TestSessioner_Register(t *testing.T) {
	ss := newMockSessioner()

	addr := "127.0.0.1:83337"

	err := ss.Register(&mockClient{RemoteAddr: addr}, &mockConn{})
	if err != nil {
		t.Logf("expect no error got %v", err)
	}
}

func TestSessioner_Unregister(t *testing.T) {
	t.Run("unregister ok", func(t *testing.T) {
		ss := newMockSessioner()

		addr := "127.0.0.1:83337"

		cl := &mockClient{RemoteAddr: addr}

		err := ss.Register(cl, &mockConn{})
		if err != nil {
			t.Fatal(err)
		}

		err = ss.Unregister(cl)
		if err != nil {
			t.Errorf("expect no unregister error got %v", err)
		}
	})

	t.Run("try to unregister non-existing client", func(t *testing.T) {
		ss := newMockSessioner()

		addr := "127.0.0.1:83337"

		cl := &mockClient{RemoteAddr: addr}

		err := ss.Register(cl, &mockConn{})
		if err != nil {
			t.Fatal(err)
		}

		anotherCl := &mockClient{RemoteAddr: "127.0.0.1:83338"}

		err = ss.Unregister(anotherCl)
		if !errors.Is(err, ErrClientNotFound) {
			t.Errorf("expect err %v got %v", ErrClientNotFound, err)
		}
	})
}

func TestSessioner_Session(t *testing.T) {
	t.Run("session found ok", func(t *testing.T) {
		ss := newMockSessioner()

		addr := "127.0.0.1:83337"

		cl := &mockClient{RemoteAddr: addr}

		err := ss.Register(cl, &mockConn{})
		if err != nil {
			t.Fatal(err)
		}

		ses, err := ss.Session(cl.Addr())
		if err != nil {
			t.Errorf("expect no unregister error got %v", err)
		}

		if ses.ClientId() != cl.Addr() {
			t.Errorf("expect client session id %s got %s", cl.Addr(), ses.ClientId())
		}
	})

	t.Run("client not found", func(t *testing.T) {
		ss := newMockSessioner()

		addr := "127.0.0.1:83337"

		cl := &mockClient{RemoteAddr: addr}

		err := ss.Register(cl, &mockConn{})
		if err != nil {
			t.Fatal(err)
		}

		anotherCl := &mockClient{RemoteAddr: "127.0.0.1:83338"}

		_, err = ss.Session(anotherCl.Addr())
		if !errors.Is(err, ErrClientNotFound) {
			t.Errorf("expect err %v got %v", ErrClientNotFound, err)
		}
	})
}
