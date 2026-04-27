package safety

import "testing"

func TestClassifyRejectsWritableCTE(t *testing.T) {
	c := Classifier{}
	_, err := c.Classify(ParsedStatement{StatementCount: 1, RootType: "SELECT", HasWritableCTE: true})
	if err != ErrWritableCTENotSupported {
		t.Fatalf("expected ErrWritableCTENotSupported, got %v", err)
	}
}

func TestClassifyInsert(t *testing.T) {
	c := Classifier{}
	class, err := c.Classify(ParsedStatement{StatementCount: 1, RootType: "INSERT"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if class != StatementInsert {
		t.Fatalf("unexpected class: %s", class)
	}
}

func TestClassifyMultiStatement(t *testing.T) {
	c := Classifier{}
	_, err := c.Classify(ParsedStatement{StatementCount: 2, RootType: "SELECT"})
	if err != ErrMultiStatement {
		t.Fatalf("expected ErrMultiStatement, got %v", err)
	}
}
