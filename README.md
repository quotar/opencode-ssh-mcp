# opencode-ssh-mcp

通过 MCP (Model Context Protocol) 协议为 opencode 提供 SSH 远程服务器管理能力。

## 项目概述

**核心功能**: 通过 SSH Config 管理多台远程服务器，提供交互式 SSH 会话给 opencode 使用  
**目标用户**: 开发者/运维人员，需要通过 opencode 管理多台远程 Linux 服务器  
**技术栈**: Go (轻量级，编译成单一可执行文件)

## 1. 架构设计

### 1.1 整体架构

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   opencode      │◄──►│  MCP Server      │◄──►│  SSH Manager    │
│                 │    │  (Go binary)     │    │  (Core Logic)   │
└─────────────────┘    └──────────────────┘    └────────┬────────┘
                                                         │
                                      ┌──────────────────┼──────────────────┐
                                      │                  │                  │
                              ┌───────▼───────┐  ┌───────▼───────┐  ┌───────▼───────┐
                              │  SSH Session  │  │  SSH Session  │  │  SSH Session  │
                              │   (Host A)    │  │   (Host B)    │  │   (Host C)    │
                              └───────┬───────┘  └───────┬───────┘  └───────┬───────┘
                                      │                  │                  │
                              ┌───────▼───────┐  ┌───────▼───────┐  ┌───────▼───────┐
                              │ Remote Host A │  │ Remote Host B │  │ Remote Host C │
                              └───────────────┘  └───────────────┘  └───────────────┘
```

### 1.2 技术选型

| 组件 | 技术选型 | 说明 |
|------|---------|------|
| **语言** | Go 1.21+ | 轻量、编译成单一 binary |
| **SSH 库** | `golang.org/x/crypto/ssh` | 标准库，无额外依赖 |
| **MCP 协议** | 手写 JSON-RPC 2.0 | 最小化依赖 |
| **传输层** | STDIO | 与 opencode 本地通信 |

## 2. 功能规格

### 2.1 核心功能

根据用户确认的需求：

| 功能 | 说明 | 优先级 |
|------|------|--------|
| **SSH Config 解析** | 读取 `~/.ssh/config` 和 `~/.ssh/config.d/*.conf` | P0 |
| **SSH 连接管理** | 建立/维护/断开 SSH 连接，支持多连接 | P0 |
| **命令执行** | 通过 SSH 执行命令并返回结果 | P0 |
| **多主机切换** | 快速在多台主机间切换 | P1 |
| **会话复用** | 保持连接复用，避免频繁建立 | P1 |
| **自动重试** | 连接失败自动重试 3 次后报错 | P1 |
| **ProxyJump 支持** | 支持多跳 SSH (ProxyJump) | P2 |

### 2.2 用户交互方式

```bash
# 在 opencode 中使用
> ssh list
Available hosts:
  - example-server (203.0.113.10)
  - test-server (203.0.113.11)

> ssh connect example-server
Connected to example-server

> df -h
Filesystem      Size  Used Avail Use% Mounted on
/dev/sda1       100G   45G   55G  45% /

> ssh connect test-server
Switched to test-server

> uptime
 10:30:45 up 123 days,  5:23,  2 users,  load average: 0.15, 0.10, 0.05
```

### 2.3 MCP Protocol 定义

```json
{
  "tools": [
    {
      "name": "ssh_list",
      "description": "List available SSH hosts from config",
      "inputSchema": { "type": "object", "properties": {} }
    },
    {
      "name": "ssh_connect",
      "description": "Connect to a remote server via SSH",
      "inputSchema": {
        "type": "object",
        "properties": {
          "host": { "type": "string", "description": "SSH config host alias or IP" }
        },
        "required": ["host"]
      }
    },
    {
      "name": "ssh_exec",
      "description": "Execute command on remote server",
      "inputSchema": {
        "type": "object",
        "properties": {
          "command": { "type": "string", "description": "Command to execute" },
          "host": { "type": "string", "description": "Target host (optional, uses active session)" }
        },
        "required": ["command"]
      }
    },
    {
      "name": "ssh_status",
      "description": "Show current SSH connection status",
      "inputSchema": { "type": "object", "properties": {} }
    },
    {
      "name": "ssh_disconnect",
      "description": "Disconnect from a remote server",
      "inputSchema": {
        "type": "object",
        "properties": {
          "host": { "type": "string", "description": "Host alias to disconnect" }
        }
      }
    }
  ]
}
```

## 3. 目录结构

```
opencode-ssh-mcp/
├── go.mod                  # Go 模块配置
├── go.sum
├── main.go                 # MCP Server 入口
├── README.md               # 使用说明
├── LICENSE                 # 许可证
├── Makefile                # 构建脚本
├── .gitignore              # Git 忽略规则
├── Dockerfile              # 容器化支持
├── ARCHITECTURE.md         # 架构设计
├── CHANGELOG.md            # 版本历史
├── SECURITY.md             # 安全政策
├── CONTRIBUTING.md         # 贡献指南
├── CODE_OF_CONDUCT.md      # 行为准则
│
├── cmd/
│   └── ssh-manager/
│       └── main.go         # CLI 入口（可选）
│
├── internal/
│   ├── mcp/
│   │   └── server.go       # MCP Server 实现
│   │
│   ├── ssh/
│   │   ├── config.go       # SSH Config 解析（支持 ~/.ssh/config.d/）
│   │   ├── connection.go   # SSH 连接管理（多连接、重试机制）
│   │   └── executor.go     # 命令执行器
│   │
│   └── session/
│       └── manager.go      # 会话管理（活跃会话、连接池）
│
└── docs/
    ├── getting-started.md  # 快速开始
    ├── architecture.md     # 架构设计
    ├── comparison.md       # CLI vs MCP 对比
    └── faq.md              # 常见问题
```

## 4. 安装和使用

### 4.1 安装

#### 前置要求
- Go 1.21+
- opencode 或其他 MCP 兼容客户端
- SSH 配置文件 (`~/.ssh/config`) 已配置好目标服务器

#### 安装方式
```bash
# 从源码安装
go install github.com/quotar/opencode-ssh-mcp@latest

# 或下载预编译二进制
wget https://github.com/quotar/opencode-ssh-mcp/releases/latest/download/opencode-ssh-mcp_linux_amd64
chmod +x opencode-ssh-mcp_linux_amd64
```

### 4.2 配置

在 opencode 配置文件中添加：

```json
{
  "mcpServers": {
    "ssh-manager": {
      "type": "local",
      "command": "opencode-ssh-mcp",
      "args": []
    }
  }
}
```

### 4.3 SSH 配置示例

SSH 配置文件 (`~/.ssh/config`) 示例：

```
Host example-server
    HostName 203.0.113.10
    Port 22
    User myuser
    IdentityFile ~/.ssh/id_rsa

Host test-server
    HostName 203.0.113.11
    Port 22
    User myuser
    IdentityFile ~/.ssh/id_rsa
```

### 4.4 使用示例

```bash
# 在 opencode 中
> ssh list
Available hosts:
  - example-server
  - test-server

> ssh connect example-server
Connected to example-server

> ssh exec df -h
Filesystem      Size  Used Avail Use% Mounted on
...

# 或者直接执行命令（使用活跃会话）
> df -h
Filesystem      Size  Used Avail Use% Mounted on
...
```

## 5. 开发

```bash
# 克隆项目
git clone https://github.com/quotar/opencode-ssh-mcp.git
cd opencode-ssh-mcp

# 构建二进制文件
make build

# 运行测试
make test
```

## 6. 资源预期

| 指标 | 预期 |
|------|------|
| **编译后大小** | ~15MB |
| **内存占用** | ~20MB |
| **启动时间** | < 100ms |
| **依赖** | 无（静态编译） |
| **适用环境** | WSL2, Linux, macOS |

## 7. 特性

### 已完成工作

1. **SSH-MCP 服务开发**
   - 完整实现了 MCP Server（Model Context Protocol）
   - 支持 SSH 配置文件读取（`~/.ssh/config`, `~/.ssh/config.d/*.conf`）
   - 实现了连接池和会话状态管理
   - 支持自动重试和 ProxyJump

2. **功能实现**
   - `ssh_list`: 列出可用主机
   - `ssh_connect`: 连接到主机  
   - `ssh_exec`: 执行命令
   - `ssh_status`: 查看状态
   - `ssh_disconnect`: 斷开连接

3. **文档**
   - 详细架构设计
   - SSH CLI 与 MCP 方式对比分析
   - MCP vs CLI 常见问题解答
   - 为什么在此场景 MCP 更适合的分析

### 核心优势

| 方面 | SSH CLI | MCP | 优势 |
|------|---------|-----|------|
| **Token 消耗** | 20 tokens/命令 | 6 tokens/命令 | MCP (少 70%) |
| **会话管理** | 无 | 有完整的会话状态 | MCP |
| **连续操作** | 每次都重复 `ssh host` | 连接后无需重复 | MCP |
| **AI 心智负担** | 高（需记住模板） | 低（会话状态） | MCP |
| **错误处理** | 基础 | 结构化 | MCP |

## 8. 许可证

MIT 许可证 - 参见 [LICENSE](LICENSE) 文件。

## 9. 贡献

参见 [CONTRIBUTING.md](CONTRIBUTING.md) 文件。

## 10. 安全

参见 [SECURITY.md](SECURITY.md) 文件。