package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/quotar/opencode-ssh-mcp/internal/mcp"
)

func main() {
	server := mcp.NewServer()

	// 设置 STDIO 传输层（标准 MCP 模式）
	stdin := bufio.NewReader(os.Stdin)
	stdout := os.Stdout

	for {
		line, err := stdin.ReadString('\n')
		if err != nil {
			break
		}

		var request mcp.Request
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			// 跳过格式错误的请求
			continue
		}

		response := server.HandleRequest(request)
		if response != nil {
			jsonResponse, _ := json.Marshal(response)
			fmt.Fprintln(stdout, string(jsonResponse))
			stdout.Sync()
		}

		// 如果请求关闭，则退出
		if request.Method == "shutdown" {
			break
		}
	}
}
