package safety

import "testing"

func TestParamsHashTypeAware(t *testing.T) {
	h1, err := ParamsHash([]any{1})
	if err != nil {
		t.Fatal(err)
	}
	h2, err := ParamsHash([]any{"1"})
	if err != nil {
		t.Fatal(err)
	}
	if h1 == h2 {
		t.Fatalf("expected different hashes for int and string params")
	}
}

func TestASTFingerprintExplicitPlaceholder(t *testing.T) {
	fp := ASTFingerprint("select 1")
	if fp != PlaceholderASTFingerprint {
		t.Fatalf("expected explicit placeholder fingerprint, got %q", fp)
	}
}
