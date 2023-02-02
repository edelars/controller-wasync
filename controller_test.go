package controller_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/edelars/controller-wasync"
	"testing"
)

type DummyError struct {
	AnyField int
}

func (de *DummyError) Error() string {
	return fmt.Sprintf("dummyerror: %d", de.AnyField)
}

type dummyCommand struct {
	DummyArg            int
	DummyOut            int
	MustFail            bool
	MustFailCustomError bool
}

type dummyHandler struct{}

func (h *dummyHandler) Exec(ctx context.Context, args *dummyCommand) error {
	if args.MustFail {
		return errors.New("handler failed")
	}
	if args.MustFailCustomError {
		return &DummyError{AnyField: 42}
	}
	args.DummyOut = args.DummyArg
	return nil
}

func (h *dummyHandler) Context() interface{} {
	return (*dummyCommand)(nil)
}

func TestDummyHandler_ShouldBeValid(t *testing.T) {
	ctrl := controller.New()

	err := ctrl.RegisterHandler(&dummyHandler{})
	if err != nil {
		t.Errorf("failed to register handler: %s", err)
	}

	arg := dummyCommand{DummyArg: 4}
	err = ctrl.Exec(context.Background(), &arg)
	if err != nil {
		t.Errorf("failed to exec command: %d", err)
	}
	if arg.DummyOut != arg.DummyArg {
		t.Errorf("command not executed")
	}
}

func TestDummyHandler_ShouldReturnProperError(t *testing.T) {
	ctrl := controller.New()

	err := ctrl.RegisterHandler(&dummyHandler{})
	if err != nil {
		t.Errorf("failed to register handler: %s", err)
	}

	arg := dummyCommand{MustFail: true}
	err = ctrl.Exec(context.Background(), &arg)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if err.Error() != "handler failed" {
		t.Errorf("got unknown error: %s", err)
	}
}

func TestDummyHandler_ShouldBypassUserDefinedErrorType(t *testing.T) {
	ctrl := controller.New()

	err := ctrl.RegisterHandler(&dummyHandler{})
	if err != nil {
		t.Errorf("failed to register handler: %s", err)
	}

	arg := dummyCommand{MustFailCustomError: true}
	err = ctrl.Exec(context.Background(), &arg)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if err.Error() != "dummyerror: 42" {
		t.Errorf("got unknown error: %s", err)
	}
}

func TestCustomErrorOnUnknwonHandler(t *testing.T) {
	ctrl := controller.New()

	arg := dummyCommand{MustFailCustomError: true}
	err := ctrl.Exec(context.Background(), &arg)
	if !errors.Is(err, &controller.ErrHandlerNotFound{}) {
		t.Errorf("unexpected error type %v", err)
	}
}

func TestDummyHandler_ShouldBeValidASync(t *testing.T) {
	ctrl := controller.New()

	err := ctrl.RegisterHandler(&dummyHandler{})
	if err != nil {
		t.Errorf("failed to register handler: %s", err)
	}

	arg := dummyCommand{DummyArg: 4}
	errCh := ctrl.ExecASync(context.Background(), &arg)

	err = <-errCh

	if err != nil {
		t.Errorf("failed to exec command: %d", err)
	}
	if arg.DummyOut != arg.DummyArg {
		t.Errorf("command not executed")
	}
}
