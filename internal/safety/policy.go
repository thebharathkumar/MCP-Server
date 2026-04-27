package safety

import (
	"fmt"
	"strings"
)

var dangerousFunctions = map[string]struct{}{
	"pg_read_file":        {},
	"pg_read_binary_file": {},
	"pg_ls_dir":           {},
	"pg_stat_file":        {},
	"lo_import":           {},
	"lo_export":           {},
}

type Policy struct {
	Classifier Classifier
}

func NewPolicy() Policy {
	return Policy{Classifier: Classifier{}}
}

func (p Policy) CheckRunQuery(stmt ParsedStatement) error {
	class, err := p.Classifier.Classify(stmt)
	if err != nil {
		return err
	}
	if class == StatementExplain {
		return ErrUseExplainTools
	}
	if class != StatementSelect {
		return ErrRunQueryOnlySelect
	}
	if hasDangerousFunction(stmt.FunctionCalls) {
		return fmt.Errorf("run_query blocked due to dangerous function usage")
	}
	return nil
}

func (p Policy) CheckPreviewWrite(stmt ParsedStatement) ([]string, error) {
	class, err := p.Classifier.Classify(stmt)
	if err != nil {
		return nil, err
	}
	if class != StatementInsert && class != StatementUpdate && class != StatementDelete {
		return nil, ErrPreviewWriteOnlyDML
	}

	warnings := make([]string, 0, 1)
	if (class == StatementUpdate || class == StatementDelete) && !stmt.HasWhereClause {
		warnings = append(warnings, fmt.Sprintf("%s has no WHERE clause", strings.ToUpper(string(class))))
	}
	return warnings, nil
}

func (p Policy) CheckExplainQuery(stmt ParsedStatement) error {
	class, err := p.Classifier.Classify(stmt)
	if err != nil {
		return err
	}
	switch class {
	case StatementSelect, StatementInsert, StatementUpdate, StatementDelete:
		return nil
	default:
		return ErrUnknownStatement
	}
}

func (p Policy) CheckExplainAnalyzeQuery(stmt ParsedStatement) error {
	class, err := p.Classifier.Classify(stmt)
	if err != nil {
		return err
	}
	if class == StatementSelect {
		return nil
	}
	if class == StatementInsert || class == StatementUpdate || class == StatementDelete {
		return ErrExplainAnalyzeWriteBlocked
	}
	return ErrUnknownStatement
}

func hasDangerousFunction(calls []string) bool {
	for _, c := range calls {
		name := strings.ToLower(c)
		if _, ok := dangerousFunctions[name]; ok {
			return true
		}
		if strings.HasPrefix(name, "dblink") {
			return true
		}
	}
	return false
}
