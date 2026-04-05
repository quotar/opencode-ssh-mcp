package ssh

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

// Connection represents an SSH connection
type Connection struct {
	client *ssh.Client
	config *SSHConfig
	id     string
	active bool
}

// ConnectionManager manages SSH connections
type ConnectionManager struct {
	connections map[string]*Connection
	configs     map[string]*SSHConfig
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*Connection),
		configs:     make(map[string]*SSHConfig),
	}
}

// SetConfigs sets the SSH configurations
func (cm *ConnectionManager) SetConfigs(configs map[string]*SSHConfig) {
	cm.configs = configs
}

// Connect connects to a host using SSH
func (cm *ConnectionManager) Connect(host string) error {
	config, exists := cm.configs[host]
	if !exists {
		return fmt.Errorf("host '%s' not found in SSH config", host)
	}

	// Check if already connected
	if conn, exists := cm.connections[host]; exists && conn.active {
		return nil // Already connected
	}

	// Create connection configuration
	address := fmt.Sprintf("%s:%s", config.HostName, config.Port)

	// Prepare SSH client configuration
	sshConfig, err := cm.prepareSSHConfig(config)
	if err != nil {
		return fmt.Errorf("failed to prepare SSH config: %w", err)
	}

	// Establish connection with retry mechanism
	var client *ssh.Client
	maxRetries := 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		client, err = ssh.Dial("tcp", address, sshConfig)
		if err == nil {
			break
		}

		if attempt == maxRetries {
			return fmt.Errorf("failed to connect to %s after %d attempts: %w", host, maxRetries, err)
		}

		// Wait with exponential backoff
		waitTime := time.Duration(attempt*attempt) * time.Second
		time.Sleep(waitTime)
	}

	// Create Connection object
	conn := &Connection{
		client: client,
		config: config,
		id:     fmt.Sprintf("%s-%d", host, time.Now().Unix()),
		active: true,
	}

	cm.connections[host] = conn
	return nil
}

// Disconnect disconnects from a host
func (cm *ConnectionManager) Disconnect(host string) error {
	conn, exists := cm.connections[host]
	if !exists {
		return fmt.Errorf("not connected to host: %s", host)
	}

	if conn.client != nil {
		conn.client.Close()
	}

	conn.active = false
	delete(cm.connections, host)
	return nil
}

// Execute executes a command on the remote host
func (cm *ConnectionManager) Execute(host, command string) (string, string, error) {
	conn, exists := cm.connections[host]
	if !exists || !conn.active {
		return "", "", fmt.Errorf("not connected to host: %s", host)
	}

	session, err := conn.client.NewSession()
	if err != nil {
		return "", "", fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	// Create pipes to capture stdout and stderr
	stdout, err := session.StdoutPipe()
	if err != nil {
		return "", "", fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		return "", "", fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := session.Start(command); err != nil {
		return "", "", fmt.Errorf("failed to start command: %w", err)
	}

	// Read stdout and stderr
	stdoutBytes := make([]byte, 0)
	stderrBytes := make([]byte, 0)

	// Read stdout
	go func() {
		stdoutBytes, _ = io.ReadAll(stdout)
	}()

	// Read stderr
	stderrBytes, _ = io.ReadAll(stderr)

	// Wait for command to complete
	err = session.Wait()

	stdoutStr := string(stdoutBytes)
	stderrStr := string(stderrBytes)

	return stdoutStr, stderrStr, err
}

// GetActiveConnections returns the list of active connections
func (cm *ConnectionManager) GetActiveConnections() []string {
	var activeHosts []string
	for host, conn := range cm.connections {
		if conn.active {
			activeHosts = append(activeHosts, host)
		}
	}
	return activeHosts
}

// prepareSSHConfig prepares the SSH client configuration
func (cm *ConnectionManager) prepareSSHConfig(config *SSHConfig) (*ssh.ClientConfig, error) {
	authMethods := []ssh.AuthMethod{}

	// Add private key authentication if identity file is specified
	if config.IdentityFile != "" {
		privateKey, err := os.ReadFile(config.IdentityFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read private key file: %w", err)
		}

		signer, err := ssh.ParsePrivateKey(privateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	// Create SSH client configuration
	sshConfig := &ssh.ClientConfig{
		User: config.User,
		Auth: authMethods,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			// Accept any host key (not ideal for production, but good for first version)
			return nil
		},
		Timeout: 30 * time.Second,
	}

	return sshConfig, nil
}
