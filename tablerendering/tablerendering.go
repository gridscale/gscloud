package tablerendering

import (
	"fmt"
	"io"
	"strings"

	"github.com/gridscale/table"
)

// RenderTable prints a table to the given io.Writer.
func RenderTable(buf io.Writer, columns []string, rows [][]string) {

	s := make([]interface{}, len(columns))
	for i, v := range columns {
		s[i] = v
	}
	tbl := table.New(s...)

	for _, row := range rows {
		s := make([]interface{}, len(row))
		for i, v := range row {
			s[i] = v
		}
		tbl.AddRow(s...)

	}

	tbl.WithWriter(buf).Print()
}
func init() {
	table.DefaultHeaderFormatter = func(format string, vals ...interface{}) string {
		return strings.ToUpper(fmt.Sprintf(format, vals...))
	}
}
