package cmdrouter

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
)

func TestBasicRouter(t *testing.T) {
	ctx := context.Background()
	var output bytes.Buffer

	executed := false

	opts := []Option{
		{
			Name: "Test Option",
			Handler: func(ctx context.Context) error {
				executed = true
				output.WriteString("Handler executed\n")
				return nil
			},
		},
	}

	router := NewCmdRouterWithSettings("Test Menu",
		WithOptions(opts...),
		WithInputOutput(strings.NewReader("1\n0\n"), &output),
	)

	router.Run(ctx)
	if !executed {
		t.Error("Handler was not executed")
	}

	if !strings.Contains(output.String(), "Test Menu") {
		t.Error("Menu title not printed")
	}
	if !strings.Contains(output.String(), "Handler executed") {
		t.Error("Handler output missing")
	}
}

func TestMiddlewareOrder(t *testing.T) {
	ctx := context.Background()
	var output bytes.Buffer

	callOrder := []string{}

	global1 := func(next Handler) Handler {
		return func(ctx context.Context) error {
			callOrder = append(callOrder, "global1")
			return next(ctx)
		}
	}
	local1 := func(next Handler) Handler {
		return func(ctx context.Context) error {
			callOrder = append(callOrder, "local1")
			return next(ctx)
		}
	}

	opt := Option{
		Name: "Test",
		Handler: func(ctx context.Context) error {
			callOrder = append(callOrder, "handler")
			return nil
		},
	}
	opt.AddMiddlewares(local1)

	router := NewCmdRouterWithSettings("Menu",
		WithOptions(opt),
		WithMiddlewares(global1),
		WithInputOutput(strings.NewReader("1\n0\n"), &output),
	)

	router.Run(ctx)

	expectedOrder := []string{"global1", "local1", "handler"}
	for i, v := range expectedOrder {
		if callOrder[i] != v {
			t.Errorf("Middleware call order wrong, expected %v got %v", expectedOrder, callOrder)
			break
		}
	}
}

type dummyPrinter struct {
	called bool
}

func (d *dummyPrinter) PrintTable(out io.Writer, headers []string, rows [][]any) {
	d.called = true
}

func TestCustomTablePrinter(t *testing.T) {
	ctx := context.Background()
	var output bytes.Buffer

	printer := &dummyPrinter{}

	opts := []Option{
		{
			Name: "Test Option",
			Handler: func(ctx context.Context) error {
				output.WriteString("Executed\n")
				return nil
			},
		},
	}

	router := NewCmdRouterWithSettings("Menu",
		WithOptions(opts...),
		WithTablePrinter(printer),
		WithInputOutput(strings.NewReader("1\n0\n"), &output),
	)

	router.Run(ctx)

	if !printer.called {
		t.Error("Custom table printer was not called")
	}
}
