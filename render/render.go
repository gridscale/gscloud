package render

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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

// AsJSON prints objects as JSON to given io.Writer.
func AsJSON(buf io.Writer, o interface{}) {
	json, err := json.Marshal(o)
	if err != nil {
		log.Fatal(err)
	}
	buf.Write(append(json, '\n'))
}

func init() {
	table.DefaultHeaderFormatter = func(format string, vals ...interface{}) string {
		return strings.ToUpper(fmt.Sprintf(format, vals...))
	}
}
