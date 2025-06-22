package main

import (
	"context"
	"fmt"

	"github.com/hahaclassic/cmdrouter"
)

func main() {
	ctx := context.Background()

	// Middleware for authentication
	authMiddleware := func(next cmdrouter.Handler) cmdrouter.Handler {
		return func(ctx context.Context) error {
			fmt.Println("[Middleware] Authenticated!")
			return next(ctx)
		}
	}

	// Login handler
	loginHandler := cmdrouter.Option{
		Name: "Login",
		Handler: func(ctx context.Context) error {
			fmt.Println("You are now logged in!")
			return nil
		},
	}

	// Profile handler
	profileHandler := cmdrouter.Option{
		Name: "View Profile",
		Handler: func(ctx context.Context) error {
			fmt.Println("Name: John Doe\nEmail: john@example.com")
			return nil
		},
	}

	// Create the command router with handlers
	router := cmdrouter.NewCmdRouter("Main Menu", loginHandler, profileHandler)
	router.AddMiddlewares(authMiddleware)

	// Start the router
	router.Run(ctx)
}
