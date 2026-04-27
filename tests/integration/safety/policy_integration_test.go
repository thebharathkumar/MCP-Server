package safety_test

import "testing"

func TestPolicyIntegrationMatrixBlockedOnParser(t *testing.T) {
	t.Skip("blocked: pg_query_go integration is required before safety matrix integration tests are meaningful")
}

func TestParserIntegrationRequired(t *testing.T) {
	t.Skip("TODO: wire real pg_query_go parser adapter before implementing handlers")
}
