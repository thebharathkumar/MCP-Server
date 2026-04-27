package session

import (
	"errors"
	"sync"
	"time"
)

// MaxPendingWrites is fixed in v0.1 to cap memory growth from runaway preview loops.
// We can make this configurable in v0.2 once usage data justifies tuning.
const MaxPendingWrites = 16

type PendingWrite struct {
	ID        string
	ExpiresAt time.Time
}

type State struct {
	WritesEnabled bool
	Pending       map[string]PendingWrite
}

type Manager struct {
	mu    sync.Mutex
	state State
}

func NewManager() *Manager {
	return &Manager{state: State{Pending: make(map[string]PendingWrite)}}
}

func (m *Manager) EnableWrites() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.state.WritesEnabled = true
}

func (m *Manager) AddPending(p PendingWrite) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cleanupLocked(time.Now())
	if len(m.state.Pending) >= MaxPendingWrites {
		return errors.New("pending write limit reached")
	}
	m.state.Pending[p.ID] = p
	return nil
}

func (m *Manager) CleanupExpired(now time.Time) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	before := len(m.state.Pending)
	m.cleanupLocked(now)
	return before - len(m.state.Pending)
}

func (m *Manager) cleanupLocked(now time.Time) {
	for id, p := range m.state.Pending {
		if !now.Before(p.ExpiresAt) {
			delete(m.state.Pending, id)
		}
	}
}
