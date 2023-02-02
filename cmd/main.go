package main

import (
	"context"
	"github.com/edelars/controller-wasync"
)

func main() {

	//Simple init controller
	ctrl := controller.New()
	if err := NewConfigController(
		ctrl,
	); err != nil {
		panic("failed to configure rbi controller")
	}

	//Example 1
	cmd := Test{}
	errCh := ctrl.ExecASync(context.Background(), &cmd)

	err := <-errCh

	if err != nil {
		panic(err)
	}

	if cmd.Out.Success {
		println("OK ASync")
	}

	//Example 2
	cmd2 := Test{}

	if err := ctrl.Exec(context.Background(), &cmd2); err != nil {
		panic(err)
	}

	if cmd2.Out.Success {
		println("OK Sync")
	}

}

func NewConfigController(
	ctrl *controller.ControllerImpl,
) (e error) {
	propogateErr := func(err error) {
		if err != nil {
			e = err
		}
	}

	propogateErr(ctrl.RegisterHandler(NewTestHandler()))
	return
}

type Test struct {
	Out struct {
		Success bool
	}
}

type TestHandler struct{}

func NewTestHandler() *TestHandler {
	return &TestHandler{}
}

func (h *TestHandler) Exec(ctx context.Context, args *Test) (err error) {
	args.Out.Success = true
	return nil
}

func (h *TestHandler) Context() interface{} {
	return (*Test)(nil)
}
