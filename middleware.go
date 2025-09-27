package cmdrouter

import (
	"context"
	"fmt"
	"log/slog"
)

// DefaultRecoverMiddleware recovers from panics in the wrapped
// handler and returns the panic as an error.
func DefaultRecoverMiddleware(next Handler) Handler {
	return func(ctx context.Context) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic: %v", r)
			}
		}()

		return next(ctx)
	}
}

// DefaultLoggerMiddleware is a middleware that logs any error
// returned by the wrapped handler using slog.Error.
func DefaultLoggerMiddleware(next Handler) Handler {
	return func(ctx context.Context) error {
		err := next(ctx)
		if err != nil {
			slog.Error("handler", "err", err)
		}
		return err
	}
}
