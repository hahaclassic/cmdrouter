package main

import (
	"context"
	"fmt"

	"github.com/hahaclassic/cmdrouter"
)

func main() {
	ctx := context.Background()

	// Global middleware
	logMiddleware := func(next cmdrouter.Handler) cmdrouter.Handler {
		return func(ctx context.Context) error {
			fmt.Println("[Middleware] Global: Logging access")
			return next(ctx)
		}
	}

	// Handler middleware
	adminCheck := func(next cmdrouter.Handler) cmdrouter.Handler {
		return func(ctx context.Context) error {
			fmt.Println("[Middleware] Handler: Admin check passed")
			return next(ctx)
		}
	}

	// handlers
	login := cmdrouter.Option{
		Name: "Login",
		Handler: func(ctx context.Context) error {
			fmt.Println("You are now logged in!")
			return nil
		},
	}

	settings := cmdrouter.Option{
		Name: "Account Settings",
		Handler: func(ctx context.Context) error {
			fmt.Println("Change your username/email/password here.")
			return nil
		},
	}

	adminPanel := cmdrouter.Option{
		Name: "Admin Panel",
		Handler: func(ctx context.Context) error {
			fmt.Println("Welcome to the admin panel.")
			return nil
		},
	}
	adminPanel.AddMiddleware(adminCheck) // add middleware for admin panel

	logHandlers := []cmdrouter.Option{
		{
			Name: "Backend logs",
			Handler: func(ctx context.Context) error {
				fmt.Println("backend logs here.")
				return nil
			},
		},
		{
			Name: "Frontend logs",
			Handler: func(ctx context.Context) error {
				fmt.Println("frontend logs here.")
				return nil
			},
		},
	}

	router := cmdrouter.NewCmdRouter("Main Menu")
	router.PathShow(true)
	router.AddMiddleware(cmdrouter.DefaultLoggerMiddleware,
		cmdrouter.DefaultRecoverMiddleware)

	devGroup := router.Group("Developer")
	devGroup.Group("Debug Logs", logHandlers...)
	devGroup.AddOptions(cmdrouter.Option{
		Name: "System Info",
		Handler: func(ctx context.Context) error {
			fmt.Println("OS: Linux\nVersion: 1.0.0")
			return nil
		}},
	)

	_ = router.Group("Settings Group", settings)
	router.AddMiddleware(logMiddleware)
	router.AddOptions(login, adminPanel)
	router.AddOptions(cmdrouter.Option{
		Name: "!!! Panic !!!",
		Handler: func(ctx context.Context) error {
			var a []int
			a[1] = 1234 // panic: runtime error: index out of range [1] with length 0
			return nil
		},
	})
	router.Run(ctx)
}
