package print

import (
	"fmt"
	"io"
	"strings"

	"github.com/liggitt/tabwriter"
	"github.com/spf13/cast"
)

const (
	tabwriterMinWidth = 6
	tabwriterWidth    = 4
	tabwriterPadding  = 3
	tabwriterPadChar  = ' '
	tabwriterFlags    = tabwriter.RememberWidths
)

// GetNewTabWriter returns a tabwriter that translates tabbed columns in input into properly aligned text.
func GetNewTabWriter(output io.Writer) *tabwriter.Writer {
	return tabwriter.NewWriter(output, tabwriterMinWidth, tabwriterWidth, tabwriterPadding, tabwriterPadChar, tabwriterFlags)
}

// TableOptions provides default flags necessary for printing.
// Given the following flag values, a printer can be requested that knows
// how to handle printing based on these values.
type TableOptions struct {
	ColumnLabels *[]string
	ShowHeaders  bool
	SortBy       *string

	// If the
	AbsoluteTimestamps bool

	// The row of the Table that is active
	ActiveRow int
}

// NewTableOptions returns options associated with
// human-readable printing, with default values set.
func NewTableOptions() *TableOptions {
	sortBy := ""
	columnLabels := []string{}

	return &TableOptions{
		ShowHeaders:        true,
		AbsoluteTimestamps: false,
		ColumnLabels:       &columnLabels,
		SortBy:             &sortBy,
	}
}

// WriteTable prints a table to the provided output
func WriteTable(columns []string, obj []map[string]interface{}, output io.Writer, options *TableOptions) error {
	if options.ShowHeaders {
		// avoid printing headers if we have no rows to display
		if len(obj) == 0 {
			return nil
		}

		first := true
		for _, column := range columns {
			if first {
				first = false
			} else {
				fmt.Fprint(output, "\t")
			}
			fmt.Fprint(output, strings.ToUpper(column))
		}
		fmt.Fprintln(output)
	}
	// print the rows
	for _, row := range obj {
		first := true
		for _, column := range columns {

			val := searchMap(row, strings.Split(column, "."))

			if first {
				first = false
			} else {
				fmt.Fprint(output, "\t")
			}
			if val != nil {
				fmt.Fprint(output, val)
			}
		}
		fmt.Fprintln(output)
	}
	return nil
}

// searchMap recursively searches for a value for path in source map.
// Returns nil if not found.
// Note: This assumes that the path entries and map keys are lower cased.
func searchMap(source map[string]interface{}, path []string) interface{} {
	if len(path) == 0 {
		return source
	}

	next, ok := source[path[0]]
	if ok {
		// Fast path
		if len(path) == 1 {
			return next
		}

		// Nested case
		switch next.(type) {
		case map[interface{}]interface{}:
			return searchMap(cast.ToStringMap(next), path[1:])
		case map[string]interface{}:
			// Type assertion is safe here since it is only reached
			// if the type of `next` is the same as the type being asserted
			return searchMap(next.(map[string]interface{}), path[1:])
		default:
			// got a value but nested key expected, return "nil" for not found
			return nil
		}
	}
	return nil
}
