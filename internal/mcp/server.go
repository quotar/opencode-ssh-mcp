package mcp

import (
	"encoding/json"
	"fmt"

	"github.com/quotar/opencode-ssh-mcp/internal/session"
	"github.com/quotar/opencode-ssh-mcp/internal/ssh"
)

// JSON-RPC 2.0 Request
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// JSON-RPC 2.0 Response
type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Tool definition for MCP
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema,omitempty"`
}

// MCP Server interface
type Server struct {
	sessionManager *session.SessionManager
	configParser   *ssh.ConfigParser
	configsLoaded  bool // Track if configs have been loaded
}

// NewServer creates a new MCP server with SSH functionality
func NewServer() *Server {
	configParser := ssh.NewConfigParser()
	connManager := ssh.NewConnectionManager()
	sessionManager := session.NewSessionManager(connManager)

	return &Server{
		sessionManager: sessionManager,
		configParser:   configParser,
		configsLoaded:  false,
	}
}

func (s *Server) HandleRequest(req Request) *Response {
	switch req.Method {
	case "initialize":
		return s.handleInitialize(req)
	case "initialized":
		return nil // No response needed
	case "shutdown":
		return s.handleShutdown(req)
	case "list_tools":
		return s.handleListTools(req)
	case "call_tool":
		return s.handleCallTool(req)
	default:
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &Error{
				Code:    -32601, // Method not found
				Message: "Method not found",
			},
		}
	}
}

func (s *Server) handleInitialize(req Request) *Response {
	// Parse SSH config on initialization
	if err := s.configParser.ParseConfigFiles(); err != nil {
		// Log error but don't fail the initialization
		fmt.Printf("Warning: failed to parse SSH configs: %v\n", err)
	} else {
		s.configsLoaded = true
		// Debug: Print loaded configs
		configs := s.configParser.GetAllConfigs()
		fmt.Printf("DEBUG: Loaded %d SSH configs\n", len(configs))
		for host := range configs {
			fmt.Printf("DEBUG: Found host: %s\n", host)
		}

		// Pass loaded configs to connection manager
		s.sessionManager.SetConfigs(configs)
	}

	// Return capabilities and tools
	result := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{
				"dynamicRegistration": true,
			},
		},
		"serverInfo": map[string]interface{}{
			"name":    "opencode-ssh-mcp",
			"version": "0.1.0",
		},
	}

	return &Response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}
}

func (s *Server) handleShutdown(req Request) *Response {
	// Graceful shutdown
	return &Response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  true,
	}
}

func (s *Server) handleListTools(req Request) *Response {
	tools := []Tool{
		{
			Name:        "ssh_list",
			Description: "List available SSH hosts from config",
		},
		{
			Name:        "ssh_connect",
			Description: "Connect to a remote server via SSH",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"host": map[string]interface{}{
						"type":        "string",
						"description": "SSH config host alias or IP",
					},
				},
				"required": []string{"host"},
			},
		},
		{
			Name:        "ssh_exec",
			Description: "Execute command on remote server",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"command": map[string]interface{}{
						"type":        "string",
						"description": "Command to execute",
					},
					"host": map[string]interface{}{
						"type":        "string",
						"description": "Target host (optional, uses active session)",
					},
				},
				"required": []string{"command"},
			},
		},
		{
			Name:        "ssh_status",
			Description: "Show current SSH connection status",
		},
		{
			Name:        "ssh_disconnect",
			Description: "Disconnect from a remote server",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"host": map[string]interface{}{
						"type":        "string",
						"description": "Host alias to disconnect",
					},
				},
			},
		},
	}

	return &Response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  tools,
	}
}

func (s *Server) handleCallTool(req Request) *Response {
	var params map[string]interface{}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &Error{
				Code:    -32700, // Parse error
				Message: "Failed to parse tool parameters",
			},
		}
	}

	// Get tool name from params (in call_tool requests)
	var toolCall struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}

	if err := json.Unmarshal(req.Params, &toolCall); err != nil {
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &Error{
				Code:    -32602, // Invalid params
				Message: "Invalid tool call format",
			},
		}
	}

	var result interface{}

	switch toolCall.Name {
	case "ssh_list":
		result = s.handleSSHList()
	case "ssh_connect":
		result = s.handleSSHConnect(toolCall.Arguments)
	case "ssh_exec":
		result = s.handleSSHExec(toolCall.Arguments)
	case "ssh_status":
		result = s.handleSSHStatus()
	case "ssh_disconnect":
		result = s.handleSSHDisconnect(toolCall.Arguments)
	default:
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &Error{
				Code:    -32601, // Method not found
				Message: fmt.Sprintf("Unknown tool: %s", toolCall.Name),
			},
		}
	}

	// Check if result is an error
	if err, isError := result.(*Error); isError {
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   err,
		}
	}

	return &Response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}
}

func (s *Server) handleSSHList() interface{} {
	configs := s.configParser.GetAllConfigs()

	var hosts []map[string]string
	for host, config := range configs {
		hostInfo := map[string]string{
			"host":     host,
			"hostname": config.HostName,
			"port":     config.Port,
			"user":     config.User,
		}
		hosts = append(hosts, hostInfo)
	}

	result := map[string]interface{}{
		"content": []interface{}{
			map[string]interface{}{
				"type": "text",
				"text": fmt.Sprintf("Found %d SSH hosts in config", len(hosts)),
			},
			map[string]interface{}{
				"type": "text",
				"text": formatHosts(hosts),
			},
		},
	}

	return result
}

func (s *Server) handleSSHConnect(arguments json.RawMessage) interface{} {
	var args struct {
		Host string `json:"host"`
	}

	if err := json.Unmarshal(arguments, &args); err != nil {
		return &Error{
			Code:    -32602,
			Message: "Invalid parameters for ssh_connect",
		}
	}

	if args.Host == "" {
		return &Error{
			Code:    -32602,
			Message: "Host parameter is required",
		}
	}

	if err := s.sessionManager.Connect(args.Host); err != nil {
		return &Error{
			Code:    -32603,
			Message: err.Error(),
		}
	}

	return map[string]interface{}{
		"content": []interface{}{
			map[string]interface{}{
				"type": "text",
				"text": fmt.Sprintf("Successfully connected to %s", args.Host),
			},
		},
	}
}

func (s *Server) handleSSHExec(arguments json.RawMessage) interface{} {
	var args struct {
		Command string `json:"command"`
		Host    string `json:"host"`
	}

	if err := json.Unmarshal(arguments, &args); err != nil {
		return &Error{
			Code:    -32602,
			Message: "Invalid parameters for ssh_exec",
		}
	}

	if args.Command == "" {
		return &Error{
			Code:    -32602,
			Message: "Command parameter is required",
		}
	}

	stdout, stderr, err := s.sessionManager.Execute(args.Command, args.Host)

	result := map[string]interface{}{
		"content": []interface{}{
			map[string]interface{}{
				"type": "text",
				"text": fmt.Sprintf("Executed command: %s", args.Command),
			},
		},
	}

	if stderr != "" {
		result["content"] = append(result["content"].([]interface{}), map[string]interface{}{
			"type": "text",
			"text": fmt.Sprintf("STDERR:\n%s", stderr),
		})
	}

	if stdout != "" {
		result["content"] = append(result["content"].([]interface{}), map[string]interface{}{
			"type": "text",
			"text": fmt.Sprintf("STDOUT:\n%s", stdout),
		})
	}

	if err != nil {
		result["content"] = append(result["content"].([]interface{}), map[string]interface{}{
			"type": "text",
			"text": fmt.Sprintf("ERROR: %v", err),
		})
	}

	return result
}

func (s *Server) handleSSHStatus() interface{} {
	activeConnections := s.sessionManager.ListActiveConnections()
	activeHost := s.sessionManager.GetActiveHost()

	content := []interface{}{
		map[string]interface{}{
			"type": "text",
			"text": fmt.Sprintf("Active connections: %d", len(activeConnections)),
		},
	}

	if len(activeConnections) > 0 {
		content = append(content, map[string]interface{}{
			"type": "text",
			"text": fmt.Sprintf("Active hosts: %v", activeConnections),
		})
	}

	if activeHost != "" {
		content = append(content, map[string]interface{}{
			"type": "text",
			"text": fmt.Sprintf("Current active host: %s", activeHost),
		})
	} else {
		content = append(content, map[string]interface{}{
			"type": "text",
			"text": "No active host selected",
		})
	}

	return map[string]interface{}{
		"content": content,
	}
}

func (s *Server) handleSSHDisconnect(arguments json.RawMessage) interface{} {
	var args struct {
		Host string `json:"host"`
	}

	if err := json.Unmarshal(arguments, &args); err != nil {
		return &Error{
			Code:    -32602,
			Message: "Invalid parameters for ssh_disconnect",
		}
	}

	host := args.Host
	if host == "" {
		// If no host specified, disconnect from active host
		host = s.sessionManager.GetActiveHost()
		if host == "" {
			return &Error{
				Code:    -32603,
				Message: "No active host to disconnect from",
			}
		}
	}

	if err := s.sessionManager.Disconnect(host); err != nil {
		return &Error{
			Code:    -32603,
			Message: err.Error(),
		}
	}

	return map[string]interface{}{
		"content": []interface{}{
			map[string]interface{}{
				"type": "text",
				"text": fmt.Sprintf("Disconnected from %s", host),
			},
		},
	}
}

// Helper function to format hosts nicely
func formatHosts(hosts []map[string]string) string {
	if len(hosts) == 0 {
		return "No SSH hosts found in config."
	}

	result := "Available SSH hosts:\n"
	for _, host := range hosts {
		result += fmt.Sprintf("  - %s (%s@%s:%s)\n", host["host"], host["user"], host["hostname"], host["port"])
	}

	return result
}
