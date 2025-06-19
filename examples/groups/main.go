package main

import (
	"context"
	"fmt"

	"github.com/hahaclassic/cmdrouter"
)

func main() {
	ctx := context.Background()

	// Global middleware
	logMiddleware := func(ctx context.Context) error {
		fmt.Println("[Middleware] Global: Logging access")
		return nil
	}

	// Handler middleware
	adminCheck := func(ctx context.Context) error {
		fmt.Println("[Middleware] Handler: Admin check passed")
		return nil
	}

	// handlers
	login := cmdrouter.OptionHandler{
		Name: "Login",
		Exec: func(ctx context.Context) error {
			fmt.Println("You are now logged in!")
			return nil
		},
	}

	settings := cmdrouter.OptionHandler{
		Name: "Account Settings",
		Exec: func(ctx context.Context) error {
			fmt.Println("Change your username/email/password here.")
			return nil
		},
	}

	adminPanel := cmdrouter.OptionHandler{
		Name: "Admin Panel",
		Exec: func(ctx context.Context) error {
			fmt.Println("Welcome to the admin panel.")
			return nil
		},
	}
	adminPanel.AddMiddleware(adminCheck) // add middleware for admin panel

	logHandlers := []cmdrouter.OptionHandler{
		{
			Name: "Backend logs",
			Exec: func(ctx context.Context) error {
				fmt.Println("backend logs here.")
				return nil
			},
		},
		{
			Name: "Frontend logs",
			Exec: func(ctx context.Context) error {
				fmt.Println("frontend logs here.")
				return nil
			},
		},
	}

	router := cmdrouter.NewCmdRouter("Main Menu", nil)
	devGroup := router.Group("Developer")
	devGroup.Group("Debug Logs", logHandlers...)
	devGroup.AddOptions(cmdrouter.OptionHandler{
		Name: "System Info",
		Exec: func(ctx context.Context) error {
			fmt.Println("OS: Linux\nVersion: 1.0.0")
			return nil
		}},
	)

	_ = router.Group("Settings Group", settings)
	router.AddMiddleware(logMiddleware)
	router.AddOptions(login, adminPanel)
	router.Run(ctx)
}
