package cmdrouter

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// TablePrinter defines the interface for printing tabular data to the console.
type TablePrinter interface {
	PrintTable(out io.Writer, headers []string, rows [][]any)
}

// Handler represents a function that processes a CLI command.
type Handler func(ctx context.Context) error

// Middleware wraps a Handler with additional logic (e.g. logging, validation, metrics).
// It takes a Handler and returns a new Handler with the middleware applied.
type Middleware func(Handler) Handler

// Option defines a CLI command with its name, execution logic, and optional middlewares.
type Option struct {
	Name        string       // Name of the operation (e.g. "login")
	Handler     Handler      // Function that executes the operation
	middlewares []Middleware // List of per-option middlewares
}

// AddMiddleware attaches a middlewares to this option.
func (o *Option) AddMiddlewares(m ...Middleware) {
	o.middlewares = append(o.middlewares, m...)
}

// Run executes the Option by wrapping its Handler with all attached middlewares in order,
// and then invoking the resulting Handler with the provided context.
// Middlewares are applied in the order they were added.
func (o *Option) Run(ctx context.Context) error {
	handler := o.Handler
	for i := len(o.middlewares) - 1; i >= 0; i-- {
		handler = o.middlewares[i](handler)
	}

	return handler(ctx)
}

// CmdRouter represents the main CLI router that handles user input and dispatches commands.
type CmdRouter struct {
	name         string       // Display name of the router or menu section.
	options      []Option     // List of available command handlers in this router.
	middlewares  []Middleware // Global middlewares applied before each handler runs.
	tablePrinter TablePrinter // Table printer used for rendering CLI menus.
	isGroup      bool         // Indicates whether this router is a subgroup (submenu).
	path         string       // Full path of this router in the CLI hierarchy, e.g. "/auth/login".
	pathShow     bool         // If true, the path is shown at the top of the menu.
	in           io.Reader    // defaults to os.Stdin
	out          io.Writer    // defaults to os.Stdout
}

// NewCmdRouter creates a new command router with the given name and optional handlers.
// It uses DefaultPrinter for printing tables and stdin/stdout for i/o streams.
func NewCmdRouter(name string, options ...Option) *CmdRouter {
	return &CmdRouter{
		name:         name,
		options:      options,
		tablePrinter: DefaultPrinter{},
		isGroup:      false,
		path:         constructPath(name),
		pathShow:     false,
		in:           os.Stdin,
		out:          os.Stdout,
	}
}

// Setting is a functional option used to configure a CmdRouter.
type Setting func(c *CmdRouter)

// NewCmdRouterWithSettings creates a new CmdRouter and applies the given settings.
func NewCmdRouterWithSettings(name string, settings ...Setting) *CmdRouter {
	router := NewCmdRouter(name)
	for _, setting := range settings {
		setting(router)
	}
	return router
}

// WithTablePrinter sets the table printer for the CmdRouter.
func WithTablePrinter(printer TablePrinter) Setting {
	return func(c *CmdRouter) {
		c.SetTablePrinter(printer)
	}
}

// WithPath enables or disables path display in the CmdRouter.
func WithPath(enable bool) Setting {
	return func(c *CmdRouter) {
		c.PathShow(enable)
	}
}

// WithMiddlewares appends the given middlewares to the CmdRouter.
func WithMiddlewares(middlewares ...Middleware) Setting {
	return func(c *CmdRouter) {
		c.AddMiddlewares(middlewares...)
	}
}

// WithOptions appends the given options (commands/handlers) to the CmdRouter.
func WithOptions(options ...Option) Setting {
	return func(c *CmdRouter) {
		c.AddOptions(options...)
	}
}

// WithInputOutput sets the input and output streams for the CmdRouter.
func WithInputOutput(in io.Reader, out io.Writer) Setting {
	return func(c *CmdRouter) {
		c.SetInputOutput(in, out)
	}
}

// Setup applies additional settings to an existing CmdRouter.
func (c *CmdRouter) Setup(settings ...Setting) {
	for _, setting := range settings {
		setting(c)
	}
}

// Group creates a submenu as a nested router and registers it as an option in the current router.
func (c *CmdRouter) Group(name string, options ...Option) *CmdRouter {
	group := &CmdRouter{
		name:         name,
		options:      options,
		tablePrinter: c.tablePrinter,
		isGroup:      true,
		path:         c.path + constructPath(name),
		pathShow:     c.pathShow,
		in:           c.in,
		out:          c.out,
	}

	c.AddOptions(Option{
		Name: name,
		Handler: func(ctx context.Context) error {
			group.Run(ctx)
			return nil
		}})

	return group
}

// SetTablePrinter sets the table printer for this router and all its groups.
func (c *CmdRouter) SetTablePrinter(printer TablePrinter) {
	c.tablePrinter = printer
}

// AddMiddlewares registers a global middlewares that will run before every option.
func (c *CmdRouter) AddMiddlewares(m ...Middleware) {
	c.middlewares = append(c.middlewares, m...)
}

// AddOptions appends new options to the router.
func (c *CmdRouter) AddOptions(options ...Option) {
	c.options = append(c.options, options...)
}

// PathShow enables or disables path display for the current router and its groups.
// When enabled, the path will be printed at the top of the menu.
func (c *CmdRouter) PathShow(enable bool) {
	c.pathShow = enable
}

func (c *CmdRouter) SetInputOutput(in io.Reader, out io.Writer) {
	c.in = in
	c.out = out
}

// Run starts the main router loop: shows the menu, processes input, applies middlewares,
// and dispatches to the selected handler.
func (c *CmdRouter) Run(ctx context.Context) {
	const exitNumber = 0
	for {
		optionNumber := c.getOptionNumber()
		if optionNumber == exitNumber {
			break
		}

		handler := c.options[optionNumber-1].Run
		for i := len(c.middlewares) - 1; i >= 0; i-- {
			handler = c.middlewares[i](handler)
		}

		fmt.Fprintln(c.out)
		_ = handler(ctx)
		fmt.Fprintln(c.out)
	}
}

// getOptionNumber displays the menu and reads the user's numeric selection from stdin.
// It keeps prompting until the input is a valid option number.
func (c CmdRouter) getOptionNumber() int {
	c.showPath()
	c.showMenu()

	scanner := bufio.NewScanner(c.in)

	for {
		fmt.Fprint(c.out, "Enter option number: ")
		if !scanner.Scan() {
			if scanner.Err() != nil {
				fmt.Fprintln(c.out, "Input error. Try again.")
				continue
			}
			break
		}

		input := strings.TrimSpace(scanner.Text())
		option, err := strconv.Atoi(input)
		if err == nil && option >= 0 && option <= len(c.options) {
			return option
		}

		fmt.Fprintln(c.out, "Invalid number. Try again.")
	}

	return 0
}

// showMenu prints the command list using the configured table printer.
func (c *CmdRouter) showMenu() {
	headers := []string{"#", c.name}
	rows := make([][]any, 0, len(c.options))

	for i := range c.options {
		rows = append(rows, []any{i + 1, c.options[i].Name})
	}

	if c.isGroup {
		rows = append(rows, []any{0, "<-Back"})
	} else {
		rows = append(rows, []any{0, "Exit"})
	}

	c.tablePrinter.PrintTable(c.out, headers, rows)
	fmt.Println()
}

// showPath prints the current router path if path display is enabled.
// Useful for nested groups to provide context on the user's location in the CLI hierarchy.
func (c *CmdRouter) showPath() {
	if c.pathShow {
		fmt.Println(c.path)
	}
}

// constructPath converts a name into a CLI path component by making it lowercase
// and replacing spaces with underscores. E.g. "User Auth" -> "/user_auth".
func constructPath(name string) string {
	return "> " + name + " "
}
