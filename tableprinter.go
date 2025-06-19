package cmdrouter

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// DefaultPrinter prints tables using simple ASCII box drawing.
//
// Example output:
// +---+----------------+
// | # |     Menu       |
// +---+----------------+
// | 1 | Login          |
// | 2 | View Profile   |
// | 0 | Exit           |
// +---+----------------+
type DefaultPrinter struct{}

// PrintTable implements the TablePrinter interface.
func (p DefaultPrinter) PrintTable(headers []string, rows [][]any) {
	if len(headers) == 0 {
		return
	}

	colWidths := p.computeColumnWidths(headers, rows)
	p.printBorder(colWidths)
	p.printRow(colWidths, p.toAny(headers))
	p.printBorder(colWidths)

	for _, row := range rows {
		p.printRow(colWidths, row)
	}

	p.printBorder(colWidths)
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
func (DefaultPrinter) printBorder(colWidths []int) {
	var b strings.Builder
	for _, w := range colWidths {
		b.WriteString("+")
		b.WriteString(strings.Repeat("-", w+2))
	}
	b.WriteString("+")
	fmt.Println(b.String())
}

// printRow prints a single row with given column widths.
func (DefaultPrinter) printRow(colWidths []int, row []any) {
	for i, cell := range row {
		format := fmt.Sprintf("| %%-%dv ", colWidths[i])
		fmt.Printf(format, cell)
	}
	fmt.Println("|")
}

// toAny converts []string to []any for uniform row printing.
func (DefaultPrinter) toAny(strs []string) []any {
	result := make([]any, len(strs))
	for i, s := range strs {
		result[i] = s
	}
	return result
}
