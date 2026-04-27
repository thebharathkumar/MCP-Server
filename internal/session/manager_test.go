package session

import (
	"fmt"
	"testing"
	"time"
)

func TestPendingLimit(t *testing.T) {
	m := NewManager()
	now := time.Now().Add(2 * time.Minute)
	for i := 0; i < MaxPendingWrites; i++ {
		if err := m.AddPending(PendingWrite{ID: fmt.Sprintf("id-%d", i), ExpiresAt: now}); err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
	}
	if err := m.AddPending(PendingWrite{ID: "overflow", ExpiresAt: now}); err == nil {
		t.Fatalf("expected limit error")
	}
}

func TestCleanupExpired(t *testing.T) {
	m := NewManager()
	now := time.Now()
	_ = m.AddPending(PendingWrite{ID: "soon-expired", ExpiresAt: now.Add(time.Second)})
	_ = m.AddPending(PendingWrite{ID: "active", ExpiresAt: now.Add(time.Minute)})

	removed := m.CleanupExpired(now.Add(2 * time.Second))
	if removed != 1 {
		t.Fatalf("expected 1 removed, got %d", removed)
	}
}
