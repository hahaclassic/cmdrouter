package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hahaclassic/cmdrouter"
	"github.com/jedib0t/go-pretty/v6/table"
)

type PrettyTablePrinter struct {
	Style table.Style
}

func (p PrettyTablePrinter) PrintTable(headers []string, rows [][]any) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
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

	authMiddleware := func(ctx context.Context) (context.Context, error) {
		fmt.Println("[Middleware] Authenticated!")
		return ctx, nil
	}

	options := []cmdrouter.OptionHandler{
		{
			Name: "Login",
			Exec: func(ctx context.Context) error {
				fmt.Println("You are now logged in!")
				return nil
			},
		},
		{
			Name: "View Profile",
			Exec: func(ctx context.Context) error {
				fmt.Println("Name: John Doe\nEmail: john@example.com")
				return nil
			},
		},
	}

	printer := PrettyTablePrinter{Style: table.StyleColoredMagentaWhiteOnBlack}
	router := cmdrouter.NewCmdRouter("Main Menu", options...)
	router.SetTablePrinter(printer)
	router.AddMiddleware(authMiddleware)

	router.Run(ctx)
}
