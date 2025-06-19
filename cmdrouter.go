package cmdrouter

import (
	"context"
	"fmt"
	"log/slog"
)

type Middleware func(ctx context.Context) error

type TablePrinter interface {
	PrintTable(headers []string, rows [][]any)
}

type OptionHandler struct {
	Name        string                          // name of the operation, e.g. "login"
	Exec        func(ctx context.Context) error // function to run the operation
	Middlewares []Middleware                    //
}

func (o *OptionHandler) AddMiddleware(m Middleware) {
	o.Middlewares = append(o.Middlewares, m)
}

func (o *OptionHandler) Run(ctx context.Context) error {
	for i := range o.Middlewares {
		if err := o.Middlewares[i](ctx); err != nil {
			return err
		}
	}

	return o.Exec(ctx)
}

type CmdRouter struct {
	name         string
	handlers     []OptionHandler
	middlewares  []Middleware
	tablePrinter TablePrinter
	isGroup      bool
}

func NewCmdRouter(name string, tablePrinter TablePrinter, handlers ...OptionHandler) *CmdRouter {
	if tablePrinter == nil {
		tablePrinter = DefaultPrinter{}
	}

	return &CmdRouter{
		name:         name,
		handlers:     handlers,
		tablePrinter: tablePrinter,
		isGroup:      false,
	}
}

func (c *CmdRouter) Group(name string, handlers ...OptionHandler) *CmdRouter {
	group := &CmdRouter{
		name:         name,
		handlers:     handlers,
		tablePrinter: c.tablePrinter,
		isGroup:      true,
	}

	c.AddOptions(OptionHandler{
		Name: name,
		Exec: func(ctx context.Context) error {
			group.Run(ctx)
			return nil
		}})

	return group
}

func (c *CmdRouter) AddMiddleware(m Middleware) {
	c.middlewares = append(c.middlewares, m)
}

func (c *CmdRouter) AddOptions(handlers ...OptionHandler) {
	c.handlers = append(c.handlers, handlers...)
}

func (c *CmdRouter) Run(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic", "err", r)
		}
	}()

	const exitNumber = 0

	for {
		option := c.getOption()
		if option == exitNumber {
			break
		}

		for i := range c.middlewares {
			if err := c.middlewares[i](ctx); err != nil {
				slog.Error("middleware", "err", err)
				continue
			}
		}

		fmt.Println()
		if err := c.handlers[option-1].Run(ctx); err != nil {
			slog.Error("handler", "err", err)
			continue
		}
		fmt.Println()
	}
}

func (c CmdRouter) getOption() int {
	c.showMenu()

	var option int
	for {
		fmt.Print("Enter option number: ")
		if _, err := fmt.Scan(&option); err == nil &&
			option >= 0 && option <= len(c.handlers) {
			break
		}

		fmt.Println("Invalid number. Try again.")
	}

	return option
}

func (c *CmdRouter) showMenu() {
	headers := []string{"#", c.name}
	rows := make([][]any, 0, len(c.handlers))

	for i := range c.handlers {
		rows = append(rows, []any{i + 1, c.handlers[i].Name})
	}

	if c.isGroup {
		rows = append(rows, []any{0, "<-Back"})
	} else {
		rows = append(rows, []any{0, "Exit"})
	}

	c.tablePrinter.PrintTable(headers, rows)
	fmt.Println()
}
