package router

import (
	"encoding/json"
	"errors"
	"github.com/Montesk/proofofwork/core/logger"
	"github.com/Montesk/proofofwork/session"
	"testing"
	"time"
)

func TestRouter_Handle(t *testing.T) {
	t.Run("handle registered route called - no error", func(t *testing.T) {
		name := "handler_name"

		rt := NewControllers(logger.NewNull())
		rt.Register(name, BuildRoute[any](func(ses session.Session, _ any) {}))

		err := rt.Handle(name, session.Null(), nil)
		if err != nil {
			t.Errorf("expect no error on route call got err %v", err)
		}
	})

	t.Run("handle method calls registered route", func(t *testing.T) {
		got := make(chan struct{})

		name := "handler_name"
		handler := func(ses session.Session, _ any) {
			got <- struct{}{}
		}

		rt := NewControllers(logger.NewNull())

		rt.Register(name, BuildRoute[any](handler))

		go rt.Handle(name, session.Null(), nil)

		for {
			select {
			case <-got:
				return // pass
			case <-time.After(1 * time.Second):
				t.Error("expect route was called")
				return
			}
		}
	})

	t.Run("route pass json message", func(t *testing.T) {
		got := make(chan struct {
			value string
		})

		msg := struct {
			value string
		}{
			value: "message in route",
		}

		name := "handler_name"
		handler := func(ses session.Session, _ any) {
			got <- msg
		}

		rt := NewControllers(logger.NewNull())

		rt.Register(name, BuildRoute[any](handler))

		raw, _ := json.Marshal(msg)

		go rt.Handle(name, session.Null(), raw)

		for {
			select {
			case gotMsg := <-got:
				if msg != gotMsg {
					t.Errorf("expect message from route %v got %v", msg, gotMsg)
					return
				}
				return // pass
			case <-time.After(1 * time.Second):
				t.Error("expect route was called")
				return
			}
		}
	})

	t.Run("handle unregistered route - expect error", func(t *testing.T) {
		name := "handler_name"

		rt := NewControllers(logger.NewNull())

		err := rt.Handle(name, session.Null(), nil)
		if err == nil {
			t.Errorf("expect no error on route call got err %v", err)
		} else {
			if !errors.Is(err, ErrRouteNotRegistered) {
				t.Errorf("expect %v err got %v", ErrRouteNotRegistered, err)
			}
		}
	})

	t.Run("another route was registered, called unregistered - expect error", func(t *testing.T) {
		name := "handler_name"

		rt := NewControllers(logger.NewNull())

		rt.Register("another_route", BuildRoute[any](func(ses session.Session, _ any) {}))

		err := rt.Handle(name, session.Null(), nil)
		if err == nil {
			t.Errorf("expect no error on route call got err %v", err)
		} else {
			if !errors.Is(err, ErrRouteNotRegistered) {
				t.Errorf("expect %v err got %v", ErrRouteNotRegistered, err)
			}
		}
	})
}
