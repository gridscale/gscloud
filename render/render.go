package render

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/gridscale/gscloud/render/table"
)

// Table prints header and rows as table to given io.Writer.
func Table(buf io.Writer, columns []string, rows [][]string) {

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
