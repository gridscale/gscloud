package render

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/gridscale/gscloud/render/table"
)

// Options holds parameters for rendering.
type Options struct {
	NoHeader bool
}

// AsTable prints header and rows as table to given io.Writer.
func AsTable(buf io.Writer, columns []string, rows [][]string, opts Options) {

	columnHeaders := make([]interface{}, len(columns))
	for i, v := range columns {
		columnHeaders[i] = v
	}
	tbl := table.New(columnHeaders...)

	for _, row := range rows {
		vals := make([]interface{}, len(row))
		for i, v := range row {
			vals[i] = v
		}
		tbl.AddRow(vals...)

	}

	tbl.WithWriter(buf).Print(!opts.NoHeader)
}

// AsJSON prints elements s JSON to given io.Writer.
func AsJSON(buf io.Writer, s ...interface{}) {
	json, _ := json.Marshal(s)
	buf.Write([]byte(fmt.Sprintf("%s\n", json)))
}

func init() {
	table.DefaultHeaderFormatter = func(format string, vals ...interface{}) string {
		return strings.ToUpper(fmt.Sprintf(format, vals...))
	}
}
