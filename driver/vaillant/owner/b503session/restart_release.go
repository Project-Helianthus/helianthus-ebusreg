package b503session

// ReleaseRestartFenceForTest clears the restart fence on m so a test can admit
// the owner gate on this same manager. It is not a recovery path.
func ReleaseRestartFenceForTest(m *Manager) {
	if m == nil {
		return
	}
	m.stateMu.Lock()
	defer m.stateMu.Unlock()
	m.restarted = false
	m.restartFences = map[byte]struct{}{}
	if m.state == Disabled && !m.mutexHeld && m.cleanup == nil {
		m.state = Idle
	}
}
