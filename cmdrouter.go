package cmdrouter

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

// Middleware represents a function that modifies the context or returns an error.
// It is typically used for logging, authentication, etc.
type Middleware func(ctx context.Context) (context.Context, error)

// TablePrinter defines the interface for printing tabular data to the console.
type TablePrinter interface {
	PrintTable(headers []string, rows [][]any)
}

// OptionHandler defines a CLI command with its name, execution logic, and optional middlewares.
type OptionHandler struct {
	Name        string                          // Name of the operation (e.g. "login")
	Exec        func(ctx context.Context) error // Function that executes the operation
	Middlewares []Middleware                    // List of per-option middlewares
}

// AddMiddleware attaches a middleware to this option.
func (o *OptionHandler) AddMiddleware(m Middleware) {
	o.Middlewares = append(o.Middlewares, m)
}

// Run executes the option by applying all middlewares and then calling the Exec function.
func (o *OptionHandler) Run(ctx context.Context) error {
	var (
		newCtx context.Context
		err    error
	)

	for _, middleware := range o.Middlewares {
		if newCtx, err = middleware(ctx); err != nil {
			return err
		}
		ctx = newCtx
	}

	return o.Exec(ctx)
}

// CmdRouter represents the main CLI router that handles user input and dispatches commands.
type CmdRouter struct {
	name         string          // Display name of the router or menu section.
	handlers     []OptionHandler // List of available command handlers in this router.
	middlewares  []Middleware    // Global middlewares applied before each handler runs.
	tablePrinter TablePrinter    // Table printer used for rendering CLI menus.
	isGroup      bool            // Indicates whether this router is a subgroup (submenu).
	path         string          // Full path of this router in the CLI hierarchy, e.g. "/auth/login".
	pathShow     bool            // If true, the path is shown at the top of the menu.
}

// NewCmdRouter creates a new command router with the given name and optional handlers.
// It uses DefaultPrinter for printing tables.
func NewCmdRouter(name string, handlers ...OptionHandler) *CmdRouter {
	return &CmdRouter{
		name:         name,
		handlers:     handlers,
		tablePrinter: DefaultPrinter{},
		isGroup:      false,
		path:         constructPath(name),
		pathShow:     false,
	}
}

// Group creates a submenu as a nested router and registers it as an option in the current router.
func (c *CmdRouter) Group(name string, handlers ...OptionHandler) *CmdRouter {
	group := &CmdRouter{
		name:         name,
		handlers:     handlers,
		tablePrinter: c.tablePrinter,
		isGroup:      true,
		path:         c.path + constructPath(name),
		pathShow:     c.pathShow,
	}

	c.AddOptions(OptionHandler{
		Name: name,
		Exec: func(ctx context.Context) error {
			group.Run(ctx)
			return nil
		}})

	return group
}

// SetTablePrinter sets the table printer for this router and all its groups.
func (c *CmdRouter) SetTablePrinter(printer TablePrinter) {
	c.tablePrinter = printer
}

// AddMiddleware registers a global middleware that will run before every option.
func (c *CmdRouter) AddMiddleware(m Middleware) {
	c.middlewares = append(c.middlewares, m)
}

// AddOptions appends new handlers to the router.
func (c *CmdRouter) AddOptions(handlers ...OptionHandler) {
	c.handlers = append(c.handlers, handlers...)
}

// PathShow enables or disables path display for the current router and its groups.
// When enabled, the path will be printed at the top of the menu.
func (c *CmdRouter) PathShow(enable bool) {
	c.pathShow = enable
}

// Run starts the main router loop: shows the menu, processes input, applies middlewares,
// and dispatches to the selected handler.
func (c *CmdRouter) Run(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic", "err", r)
		}
	}()

	const exitNumber = 0

	for {
		var (
			optionCtx context.Context = ctx
			newCtx    context.Context
			err       error
		)

		option := c.getOption()
		if option == exitNumber {
			break
		}

		isFailed := false
		for _, middleware := range c.middlewares {
			if newCtx, err = middleware(optionCtx); err != nil {
				slog.Error("middleware", "err", err)
				isFailed = true
				break
			}
			optionCtx = newCtx
		}
		if isFailed {
			continue
		}

		fmt.Println()
		if err := c.handlers[option-1].Run(optionCtx); err != nil {
			slog.Error("handler", "err", err)
			continue
		}
		fmt.Println()
	}
}

// getOption shows the menu and reads the user's selection via safe input.
func (c CmdRouter) getOption() int {
	c.showPath()
	c.showMenu()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter option number: ")
		if !scanner.Scan() {
			fmt.Println("Input error. Try again.")
			continue
		}

		input := strings.TrimSpace(scanner.Text())
		option, err := strconv.Atoi(input)
		if err == nil && option >= 0 && option <= len(c.handlers) {
			return option
		}

		fmt.Println("Invalid number. Try again.")
	}
}

// showMenu prints the command list using the configured table printer.
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

// showPath prints the current router path if path display is enabled.
// Useful for nested groups to provide context on the user's location in the CLI hierarchy.
func (c *CmdRouter) showPath() {
	if c.pathShow {
		fmt.Println(c.path)
	}
}

func constructPath(name string) string {
	return "/" + strings.ReplaceAll(strings.ToLower(name), " ", "_")
}
