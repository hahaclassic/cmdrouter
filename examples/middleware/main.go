package main

import (
	"context"
	"fmt"

	"github.com/hahaclassic/cmdrouter"
)

func main() {
	ctx := context.Background()

	// Middlewares
	global1 := func(next cmdrouter.Handler) cmdrouter.Handler {
		return func(ctx context.Context) error {
			fmt.Print("global 1 -> ")
			return next(ctx)
		}
	}

	global2 := func(next cmdrouter.Handler) cmdrouter.Handler {
		return func(ctx context.Context) error {
			fmt.Print("global 2 -> ")
			return next(ctx)
		}
	}

	global3 := func(next cmdrouter.Handler) cmdrouter.Handler {
		return func(ctx context.Context) error {
			fmt.Print("global 3 -> ")
			return next(ctx)
		}
	}

	local1 := func(next cmdrouter.Handler) cmdrouter.Handler {
		return func(ctx context.Context) error {
			fmt.Print("local 1 -> ")
			return next(ctx)
		}
	}

	local2 := func(next cmdrouter.Handler) cmdrouter.Handler {
		return func(ctx context.Context) error {
			fmt.Print("local 2 -> ")
			return next(ctx)
		}
	}

	// Handler
	handler := cmdrouter.Option{
		Name: "Handler",
		Handler: func(ctx context.Context) error {
			fmt.Println("Option Handler")
			return nil
		},
	}
	handler.AddMiddlewares(local1, local2)

	// Create the command router with handler
	router := cmdrouter.NewCmdRouter("Main Menu", handler)
	router.AddMiddlewares(
		global1,
		global2,
		global3,
	)

	// Start the router
	router.Run(ctx)
	// Result: global 1 -> global 2 -> global 3 -> local 1 -> local 2 -> Option Handler
}
