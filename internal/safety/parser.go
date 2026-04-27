package safety

import "errors"

var ErrParserNotImplemented = errors.New("stop-the-line: pg_query_go parser integration is required before policy enforcement can be trusted")

type Parser interface {
	Parse(sql string) (ParsedStatement, error)
}
