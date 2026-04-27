package safety

import "testing"

func TestRunQueryRejectsExplain(t *testing.T) {
	p := NewPolicy()
	err := p.CheckRunQuery(ParsedStatement{StatementCount: 1, RootType: "EXPLAIN"})
	if err != ErrUseExplainTools {
		t.Fatalf("expected ErrUseExplainTools, got %v", err)
	}
}

func TestRunQueryRejectsDangerousFunction(t *testing.T) {
	p := NewPolicy()
	err := p.CheckRunQuery(ParsedStatement{StatementCount: 1, RootType: "SELECT", FunctionCalls: []string{"pg_read_file"}})
	if err == nil {
		t.Fatalf("expected dangerous function rejection")
	}
}

func TestPreviewWriteWarnsNoWhere(t *testing.T) {
	p := NewPolicy()
	warnings, err := p.CheckPreviewWrite(ParsedStatement{StatementCount: 1, RootType: "UPDATE", HasWhereClause: false})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(warnings) != 1 {
		t.Fatalf("expected one warning, got %d", len(warnings))
	}
}

func TestExplainAnalyzeBlocksWrites(t *testing.T) {
	p := NewPolicy()
	err := p.CheckExplainAnalyzeQuery(ParsedStatement{StatementCount: 1, RootType: "DELETE"})
	if err != ErrExplainAnalyzeWriteBlocked {
		t.Fatalf("expected ErrExplainAnalyzeWriteBlocked, got %v", err)
	}
}
