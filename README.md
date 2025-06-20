# cmdrouter

[![Go Reference](https://pkg.go.dev/badge/github.com/hahaclassic/cmdrouter.svg)](https://pkg.go.dev/github.com/hahaclassic/cmdrouter)
[![Go Report Card](https://goreportcard.com/badge/github.com/hahaclassic/cmdrouter)](https://goreportcard.com/report/github.com/hahaclassic/cmdrouter)

<!-- [![Build Status](https://github.com/hahaclassic/go-pretty/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/hahaclassic/go-pretty/actions?query=workflow%3ACI+event%3Apush+branch%3Amain)
[![Coverage Status](https://coveralls.io/repos/github/hahaclassic/go-pretty/badge.svg?branch=main)](https://coveralls.io/github/hahaclassic/go-pretty?branch=main) -->

`cmdrouter` is a lightweight zero-dependency Go package for building command-line menus.

## Features

- Simple ASCII table menu printing by default
- Support for middleware that run before each command
- Grouping of commands into submenus
- No external dependencies (only Go standard library)
- Customizable table output by implementing the `TablePrinter` interface
- Optionally, you can use libraries like [`go-pretty`](https://github.com/jedib0t/go-pretty) for prettier tables

## Install

```bash
    go get github.com/hahaclassic/cmdrouter
```

## Examples

More examples are [here](/examples/).

### Default

```go
package main

import (
    "context"
    "fmt"

    "github.com/hahaclassic/cmdrouter"
)

func main() {
    ctx := context.Background()

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

    router := cmdrouter.NewCmdRouter("Main Menu", options...)
    router.Run(ctx)
}
```

Result:
```
+---+--------------+
| # | Main Menu    |
+---+--------------+
| 1 | Login        |
| 2 | View Profile |
| 0 | Exit         |
+---+--------------+

Enter option number: 2

Name: John Doe
Email: john@example.com
```

## Groups

Groups allow nesting commands under a submenu to better organize related options.
Each group is itself a CmdRouter with its own set of handlers and shares the same TablePrinter.

```go
func (c *CmdRouter) Group(name string, handlers ...OptionHandler) *CmdRouter
```

- Use ```router.Group()``` to define a submenu.
- Selecting the group in the CLI opens its submenu.
- `0 <-Back` is added automatically to return to the previous level.

Example:

```go
// ...
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

router := cmdrouter.NewCmdRouter("Main Menu")
devGroup := router.Group("Developer")
devGroup.Group("Debug Logs", logHandlers...)
// ...
```

Result:
```
+---+-------------+
| # | Developer   |
+---+-------------+
| 1 | Debug Logs  |
| 2 | System Info |
| 0 | <-Back      |
+---+-------------+

Enter option number: 1

+---+---------------+
| # | Debug Logs    |
+---+---------------+
| 1 | Backend logs  |
| 2 | Frontend logs |
| 0 | <-Back        |
+---+---------------+

Enter option number: 1

backend logs here.
```

## Middlewares

Use middlewares to:
- Inject values into the context
- Handle authentication or logging
- Or for any other custom processing

```go
type Middleware func(ctx context.Context) (context.Context, error)
```

### There are two types of middlewares:
- Global: Added to the router via AddMiddleware, applied to all handlers.
- Local: Added to individual handlers via OptionHandler.AddMiddleware.

### Execution Order
Middlewares are executed in the order they are added:

1. Router-level (global) middlewares

2. Handler-level (local) middlewares

3. Command execution (Exec)

Example:
```go
router.AddMiddleware(func(ctx context.Context) (context.Context, error) {
		fmt.Println("[Middleware] Authenticated!")
		return ctx, nil
	})
```

## Custom table printing

By default, cmdrouter uses a simple ASCII printer (DefaultPrinter) relying only on Go's standard library.

If you want a prettier table output, you can implement the TablePrinter interface yourself. For example, using [`go-pretty`](https://github.com/jedib0t/go-pretty):

```go
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
    // ...
    printer := PrettyTablePrinter{Style: table.StyleRounded}
	router := cmdrouter.NewCmdRouter("Main Menu", options...)
    router.SetTablePrinter(printer) // set pretty printer
    router.Run(ctx)
}
```

Result (table.StyleRounded):
```
╭───┬──────────────╮
│ # │ MAIN MENU    │
├───┼──────────────┤
│ 1 │ Login        │
│ 2 │ View Profile │
│ 0 │ Exit         │
╰───┴──────────────╯
```

You also can use table.StyleColoredMagentaWhiteOnBlack or others.

## License

Licensed under [MIT License](./LICENSE).