package ssh

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SSHConfig represents SSH configuration for a host
type SSHConfig struct {
	Host         string
	HostName     string
	Port         string
	User         string
	IdentityFile string
	ProxyJump    string
}

// ConfigParser parses SSH config files
type ConfigParser struct {
	configs map[string]*SSHConfig
}

// NewConfigParser creates a new SSH config parser
func NewConfigParser() *ConfigParser {
	return &ConfigParser{
		configs: make(map[string]*SSHConfig),
	}
}

// ParseConfigFiles parses ~/.ssh/config and ~/.ssh/config.d/*.conf
func (p *ConfigParser) ParseConfigFiles() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Parse main config file
	mainConfigPath := filepath.Join(homeDir, ".ssh", "config")
	if err := p.parseConfigFile(mainConfigPath); err != nil {
		// If main config doesn't exist, that's okay
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to parse main config: %w", err)
		}
	}

	// Parse config.d directory
	configDir := filepath.Join(homeDir, ".ssh", "config.d")
	files, err := os.ReadDir(configDir)
	if err != nil {
		// If config.d doesn't exist, that's okay
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to read config.d directory: %w", err)
		}
	} else {
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".conf") {
				configPath := filepath.Join(configDir, file.Name())
				if err := p.parseConfigFile(configPath); err != nil {
					return fmt.Errorf("failed to parse config file %s: %w", configPath, err)
				}
			}
		}
	}

	return nil
}

// parseConfigFile parses a single SSH config file
func (p *ConfigParser) parseConfigFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentHost := ""

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		key := strings.ToLower(parts[0])
		value := strings.Join(parts[1:], " ")

		switch key {
		case "host":
			currentHost = value
			// Initialize config if not exists
			if _, exists := p.configs[currentHost]; !exists {
				p.configs[currentHost] = &SSHConfig{
					Host: currentHost,
					Port: "22", // Default SSH port
				}
			}
		case "hostname":
			if currentHost != "" {
				p.configs[currentHost].HostName = value
			}
		case "port":
			if currentHost != "" {
				p.configs[currentHost].Port = value
			}
		case "user":
			if currentHost != "" {
				p.configs[currentHost].User = value
			}
		case "identityfile":
			if currentHost != "" {
				p.configs[currentHost].IdentityFile = expandPath(value)
			}
		case "proxyjump":
			if currentHost != "" {
				p.configs[currentHost].ProxyJump = value
			}
		}
	}

	return scanner.Err()
}

// GetAllConfigs returns all parsed SSH configurations
func (p *ConfigParser) GetAllConfigs() map[string]*SSHConfig {
	return p.configs
}

// GetConfig returns SSH configuration for a specific host
func (p *ConfigParser) GetConfig(host string) (*SSHConfig, bool) {
	config, exists := p.configs[host]
	return config, exists
}

// expandPath expands ~ to home directory
func expandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(homeDir, path[1:])
	}
	return path
}
