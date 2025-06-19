package main

import (
	"context"
	"fmt"

	"github.com/hahaclassic/cmdrouter"
)

func main() {
	ctx := context.Background()

	// Middleware for authentication
	authMiddleware := func(ctx context.Context) (context.Context, error) {
		fmt.Println("[Middleware] Authenticated!")
		return ctx, nil
	}

	// Login handler
	loginHandler := cmdrouter.OptionHandler{
		Name: "Login",
		Exec: func(ctx context.Context) error {
			fmt.Println("You are now logged in!")
			return nil
		},
	}

	// Profile handler
	profileHandler := cmdrouter.OptionHandler{
		Name: "View Profile",
		Exec: func(ctx context.Context) error {
			fmt.Println("Name: John Doe\nEmail: john@example.com")
			return nil
		},
	}

	// Create the command router with handlers
	router := cmdrouter.NewCmdRouter("Main Menu", loginHandler, profileHandler)
	router.AddMiddleware(authMiddleware)

	// Start the router
	router.Run(ctx)
}
