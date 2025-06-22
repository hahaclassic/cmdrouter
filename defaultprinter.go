package cmdrouter

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

// DefaultPrinter prints tables using simple ASCII box drawing.
//
//	+---+----------------+
//	| # |     Menu       |
//	+---+----------------+
//	| 1 | Login          |
//	| 2 | View Profile   |
//	| 0 | Exit           |
//	+---+----------------+
type DefaultPrinter struct{}

// PrintTable implements the TablePrinter interface.
func (p DefaultPrinter) PrintTable(out io.Writer, headers []string, rows [][]any) {
	if len(headers) == 0 {
		return
	}

	colWidths := p.computeColumnWidths(headers, rows)
	p.printBorder(out, colWidths)
	p.printRow(out, colWidths, p.toAny(headers))
	p.printBorder(out, colWidths)

	for _, row := range rows {
		p.printRow(out, colWidths, row)
	}

	p.printBorder(out, colWidths)
}

// computeColumnWidths calculates the maximum width for each column based on headers and data.
func (DefaultPrinter) computeColumnWidths(headers []string, rows [][]any) []int {
	colWidths := make([]int, len(headers))
	for i, h := range headers {
		colWidths[i] = utf8.RuneCountInString(h)
	}

	for _, row := range rows {
		for i, cell := range row {
			length := utf8.RuneCountInString(fmt.Sprint(cell))
			if length > colWidths[i] {
				colWidths[i] = length
			}
		}
	}

	return colWidths
}

// printBorder prints the horizontal border line based on column widths.
func (DefaultPrinter) printBorder(out io.Writer, colWidths []int) {
	const offset = 2
	var border strings.Builder

	for _, w := range colWidths {
		border.WriteByte('+')
		border.WriteString(strings.Repeat("-", w+offset))
	}
	border.WriteByte('+')

	_, _ = fmt.Fprintln(out, border.String())
}

// printRow prints a single row with given column widths.
func (DefaultPrinter) printRow(out io.Writer, colWidths []int, row []any) {
	for i, cell := range row {
		format := fmt.Sprintf("| %%-%dv ", colWidths[i])
		_, _ = fmt.Fprintf(out, format, cell)
	}
	_, _ = fmt.Fprintln(out, "|")
}

// toAny converts []string to []any for uniform row printing.
func (DefaultPrinter) toAny(strs []string) []any {
	result := make([]any, len(strs))
	for i, s := range strs {
		result[i] = s
	}
	return result
}
