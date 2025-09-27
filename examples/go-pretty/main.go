package main

import (
	"context"
	"fmt"
	"io"

	"github.com/hahaclassic/cmdrouter"
	"github.com/jedib0t/go-pretty/v6/table"
)

// NOTE: it is v1.0.0 version

type PrettyTablePrinter struct {
	Style table.Style
}

func (p PrettyTablePrinter) PrintTable(out io.Writer,
	headers []string, rows [][]any) {

	t := table.NewWriter()
	t.SetOutputMirror(out)
	t.SetStyle(p.Style)

	// Convert headers to table.Row
	headerRow := make(table.Row, len(headers))
	for i, h := range headers {
		headerRow[i] = h
	}
	t.AppendHeader(headerRow)

	// Append data rows
	for _, row := range rows {
		t.AppendRow(row)
	}

	t.Render()
}

func main() {
	ctx := context.Background()

	authMiddleware := func(next cmdrouter.Handler) cmdrouter.Handler {
		return func(ctx context.Context) error {
			fmt.Println("[Middleware] Authenticated!")
			return next(ctx)
		}
	}

	options := []cmdrouter.Option{
		{
			Name: "Login",
			Handler: func(ctx context.Context) error {
				fmt.Println("You are now logged in!")
				return nil
			},
		},
		{
			Name: "View Profile",
			Handler: func(ctx context.Context) error {
				fmt.Println("Name: John Doe\nEmail: john@example.com")
				return nil
			},
		},
	}

	printer := PrettyTablePrinter{Style: table.StyleColoredMagentaWhiteOnBlack}
	router := cmdrouter.NewCmdRouterWithSettings("Main Menu",
		cmdrouter.WithOptions(options...),
		cmdrouter.WithTablePrinter(printer),
		cmdrouter.WithMiddlewares(authMiddleware),
	)

	router.Run(ctx)
}
