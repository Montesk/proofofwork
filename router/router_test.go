package router

import (
	"encoding/json"
	"errors"
	"github.com/Montesk/proofofwork/handlers"
	"github.com/Montesk/proofofwork/session"
	"testing"
	"time"
)

func TestRouter_Handle(t *testing.T) {
	t.Run("handle registered handler called - no error", func(t *testing.T) {
		name := "handler_name"

		rt := New()
		rt.Register(name, handlers.BuildRoute[any](func(ses session.Session, _ any) {}))

		err := rt.Handle(name, session.Null(), nil)
		if err != nil {
			t.Errorf("expect no error on handler call got err %v", err)
		}
	})

	t.Run("handle method calls registered handler", func(t *testing.T) {
		got := make(chan struct{})

		name := "handler_name"
		handler := func(ses session.Session, _ any) {
			got <- struct{}{}
		}

		rt := New()

		rt.Register(name, handlers.BuildRoute[any](handler))

		go rt.Handle(name, session.Null(), nil)

		for {
			select {
			case <-got:
				return // pass
			case <-time.After(1 * time.Second):
				t.Error("expect handler was called")
				return
			}
		}
	})

	t.Run("handler pass json message", func(t *testing.T) {
		got := make(chan struct {
			value string
		})

		msg := struct {
			value string
		}{
			value: "message in handler",
		}

		name := "handler_name"
		handler := func(ses session.Session, _ any) {
			got <- msg
		}

		rt := New()

		rt.Register(name, handlers.BuildRoute[any](handler))

		raw, _ := json.Marshal(msg)

		go rt.Handle(name, session.Null(), raw)

		for {
			select {
			case gotMsg := <-got:
				if msg != gotMsg {
					t.Errorf("expect message from handler %v got %v", msg, gotMsg)
					return
				}
				return // pass
			case <-time.After(1 * time.Second):
				t.Error("expect handler was called")
				return
			}
		}
	})

	t.Run("handle unregistered handler - expect error", func(t *testing.T) {
		name := "handler_name"

		rt := New()

		err := rt.Handle(name, session.Null(), nil)
		if err == nil {
			t.Errorf("expect no error on handler call got err %v", err)
		} else {
			if !errors.Is(err, ErrRouteNotRegistered) {
				t.Errorf("expect %v err got %v", ErrRouteNotRegistered, err)
			}
		}
	})

	t.Run("another handler was registered, called unregistered - expect error", func(t *testing.T) {
		name := "handler_name"

		rt := New()

		rt.Register("another_route", handlers.BuildRoute[any](func(ses session.Session, _ any) {}))

		err := rt.Handle(name, session.Null(), nil)
		if err == nil {
			t.Errorf("expect no error on handler call got err %v", err)
		} else {
			if !errors.Is(err, ErrRouteNotRegistered) {
				t.Errorf("expect %v err got %v", ErrRouteNotRegistered, err)
			}
		}
	})
}
