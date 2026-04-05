package session

import (
	"sync"

	"github.com/quotar/opencode-ssh-mcp/internal/ssh"
)

// SessionManager manages SSH sessions and active hosts
type SessionManager struct {
	connectionManager *ssh.ConnectionManager
	activeHost        string
	mu                sync.RWMutex
}

// NewSessionManager creates a new session manager
func NewSessionManager(cm *ssh.ConnectionManager) *SessionManager {
	return &SessionManager{
		connectionManager: cm,
		activeHost:        "",
	}
}

// SetConfigs passes SSH configurations to the underlying connection manager
func (sm *SessionManager) SetConfigs(configs map[string]*ssh.SSHConfig) {
	sm.connectionManager.SetConfigs(configs)
}

// Connect connects to a host and makes it active
func (sm *SessionManager) Connect(host string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if err := sm.connectionManager.Connect(host); err != nil {
		return err
	}

	sm.activeHost = host
	return nil
}

// Disconnect disconnects from a host
func (sm *SessionManager) Disconnect(host string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	err := sm.connectionManager.Disconnect(host)

	// If disconnected host was the active one, clear active host
	if sm.activeHost == host {
		// Find another active connection or clear active host
		activeConnections := sm.connectionManager.GetActiveConnections()
		if len(activeConnections) > 0 {
			sm.activeHost = activeConnections[0] // Set first active connection as active
		} else {
			sm.activeHost = ""
		}
	}

	return err
}

// Execute executes a command on a host (uses active host if none specified)
func (sm *SessionManager) Execute(command string, host string) (string, string, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	targetHost := host
	if targetHost == "" {
		targetHost = sm.activeHost
		if targetHost == "" {
			return "", "", &NoActiveHostException{}
		}
	}

	return sm.connectionManager.Execute(targetHost, command)
}

// ListActiveConnections returns the list of active connections
func (sm *SessionManager) ListActiveConnections() []string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return sm.connectionManager.GetActiveConnections()
}

// GetActiveHost returns the currently active host
func (sm *SessionManager) GetActiveHost() string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return sm.activeHost
}

// NoActiveHostException indicates no active connection exists
type NoActiveHostException struct{}

func (e *NoActiveHostException) Error() string {
	return "no active SSH connection, use ssh_connect first"
}
