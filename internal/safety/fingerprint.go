package safety

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

const PlaceholderASTFingerprint = "PLACEHOLDER_FINGERPRINT_NOT_PRODUCTION_SAFE"

func ParamsHash(params []any) (string, error) {
	b, err := json.Marshal(params)
	if err != nil {
		return "", err
	}
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:]), nil
}

// ASTFingerprint is intentionally a placeholder until pg_query_go integration is available.
// Do not use this value to gate production commit semantics.
func ASTFingerprint(_ string) string {
	return PlaceholderASTFingerprint
}
