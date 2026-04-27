package safety

import "errors"

type StatementClass string

const (
	StatementUnknown StatementClass = "unknown"
	StatementSelect  StatementClass = "select"
	StatementInsert  StatementClass = "insert"
	StatementUpdate  StatementClass = "update"
	StatementDelete  StatementClass = "delete"
	StatementExplain StatementClass = "explain"
	StatementDDL     StatementClass = "ddl"
	StatementAdmin   StatementClass = "admin"
	StatementLock    StatementClass = "lock"
	StatementListen  StatementClass = "listen"
	StatementSetRole StatementClass = "set_role"
)

type ParsedStatement struct {
	RootType            string
	StatementCount      int
	HasWritableCTE      bool
	FunctionCalls       []string
	HasWhereClause      bool
	ContainsExplain     bool
	ExplainAnalyze      bool
	ExplainedRootType   string
}

var (
	ErrMultiStatement             = errors.New("only a single SQL statement is allowed")
	ErrUnknownStatement           = errors.New("unsupported or unknown SQL statement")
	ErrWritableCTENotSupported    = errors.New("data-modifying CTEs are not supported in v0.1")
	ErrRunQueryOnlySelect         = errors.New("run_query only supports read-only SELECT statements")
	ErrUseExplainTools            = errors.New("EXPLAIN statements are not allowed in run_query; use explain_query or explain_analyze_query")
	ErrPreviewWriteOnlyDML        = errors.New("preview_write supports only top-level INSERT, UPDATE, DELETE")
	ErrExplainAnalyzeWriteBlocked = errors.New("explain_analyze_query on write statements is blocked in v0.1")
)
