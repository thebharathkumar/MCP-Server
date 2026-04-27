package safety

import "strings"

type Classifier struct{}

func (Classifier) Classify(stmt ParsedStatement) (StatementClass, error) {
	if stmt.StatementCount != 1 {
		return StatementUnknown, ErrMultiStatement
	}
	if stmt.HasWritableCTE {
		return StatementUnknown, ErrWritableCTENotSupported
	}

	t := strings.ToUpper(stmt.RootType)
	switch t {
	case "SELECT":
		return StatementSelect, nil
	case "INSERT":
		return StatementInsert, nil
	case "UPDATE":
		return StatementUpdate, nil
	case "DELETE":
		return StatementDelete, nil
	case "EXPLAIN":
		return StatementExplain, nil
	case "SETROLE", "SETSESSIONAUTH":
		return StatementSetRole, nil
	case "LOCK":
		return StatementLock, nil
	case "LISTEN", "UNLISTEN":
		return StatementListen, nil
	case "CREATE", "ALTER", "DROP", "TRUNCATE":
		return StatementDDL, nil
	case "VACUUM", "ANALYZE", "REINDEX", "CLUSTER", "COPY", "DO":
		return StatementAdmin, nil
	default:
		return StatementUnknown, ErrUnknownStatement
	}
}
