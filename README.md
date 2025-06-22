# cmdrouter

[![Go Reference](https://pkg.go.dev/badge/github.com/hahaclassic/cmdrouter.svg)](https://pkg.go.dev/github.com/hahaclassic/cmdrouter)
[![Go Report Card](https://goreportcard.com/badge/github.com/hahaclassic/cmdrouter)](https://goreportcard.com/report/github.com/hahaclassic/cmdrouter)

<!-- [![Build Status](https://github.com/hahaclassic/go-pretty/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/hahaclassic/go-pretty/actions?query=workflow%3ACI+event%3Apush+branch%3Amain)
[![Coverage Status](https://coveralls.io/repos/github/hahaclassic/go-pretty/badge.svg?branch=main)](https://coveralls.io/github/hahaclassic/go-pretty?branch=main) -->

`cmdrouter` is a lightweight zero-dependency Go package for building command-line menus.

## Features

- Simple ASCII table menu printing by default
- Support for global/local middlewares
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
func (c *CmdRouter) Group(name string, handlers ...Options) *CmdRouter
```

- Use ```router.Group()``` to define a submenu.
- Selecting the group in the CLI opens its submenu.
- `0 <-Back` is added automatically to return to the previous level.

[Example](./examples/groups/main.go):

```go
// ...
logHandlers := []cmdrouter.Options{
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
router.PathShow(true) // You also can 
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
type Handler func(ctx context.Context) error

type Middleware func(Handler) Handler
```

### There are two types of middlewares:
- Global: Added to the router via AddMiddlewares, applied to all handlers.
- Local: Added to individual handlers via OptionHandler.AddMiddlewares.

```go
router.AddMiddlewares(func(next cmdrouter.Handler) cmdrouter.Handler {
		return func(ctx context.Context) error {
            fmt.Println("[Middleware] Before!")
		    err := next(ctx)
            fmt.Println("[Middleware] After!")

            return err
        }
	})
```

### Execution Order
Middlewares are executed in the order they are added:

1. Router-level (global) middlewares

2. Handler-level (local) middlewares

3. Command execution (Handler)

[Example](./examples/middleware/main.go):
```go
// ...
handler := cmdrouter.Option{
    Name: "Handler",
    Handler: func(ctx context.Context) error {
        fmt.Println("Option Handler")
        return nil
    },
}
handler.AddMiddlewares(local1, local2) // add local middlewares for this handler

router := cmdrouter.NewCmdRouter("Main Menu", handler)
router.AddMiddlewares( // add global middlewares for this router
    global1,
    global2,
    global3,
)
router.Run(ctx)
// ...
```
Result:
```
global 1 -> global 2 -> global 3 -> local 1 -> local 2 -> Option Handler
```

## Custom table printing

By default, cmdrouter uses a simple ASCII printer (DefaultPrinter) relying only on Go's standard library.

If you want a prettier table output, you can implement the TablePrinter interface yourself. For example, using [`go-pretty`](https://github.com/jedib0t/go-pretty):

```go
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

## Other features

### Path display

CmdRouter can show the current menu path to indicate nested command locations, e.g.:

```
/main_menu/developer/debug_logs
+---+---------------+
| # | Debug Logs    |
+---+---------------+
| 1 | Backend logs  |
| 2 | Frontend logs |
| 0 | <-Back        |
+---+---------------+
```

To enable path display, use the method:
```go
router.PathShow(true)
```
or the functional option ```WithPath(true)``` when creating or configuring the router.

### Settings (functional options)
CmdRouter supports flexible configuration via functional options called Settings. This allows you to conveniently customize your router with various options such as custom table printers, middlewares, path display, input/output streams, and commands.

Example of creating a router with settings:
```go
router := cmdrouter.NewCmdRouterWithSettings("Main Menu",
    cmdrouter.WithPath(true),
    cmdrouter.WithTablePrinter(myCustomPrinter),
    cmdrouter.WithMiddlewares(myMiddleware),
    cmdrouter.WithOptions(myOptions...),
)
```
Or applying settings to an existing router:

```go
router.Setup(
    cmdrouter.WithPath(true),
    cmdrouter.WithMiddlewares(additionalMiddleware),
)
```

#### Available settings include:

- WithTablePrinter(TablePrinter) — set a custom table printer

- WithPath(bool) — enable or disable path display

- WithMiddlewares(...Middleware) — add global middlewares

- WithOptions(...Option) — add command options

- WithInputOutput(io.Reader, io.Writer) — specify custom input/output streams (useful for testing, etc.)

## License

Licensed under [MIT License](./LICENSE).