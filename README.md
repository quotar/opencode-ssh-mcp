# opencode-ssh-mcp

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![GitHub release](https://img.shields.io/github/release/quotar/opencode-ssh-mcp.svg)](https://github.com/quotar/opencode-ssh-mcp/releases)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org)

通过 MCP (Model Context Protocol) 协议为 opencode 提供 SSH 远程服务器管理能力。

## 项目概述

**核心功能**: 通过 SSH Config 管理多台远程服务器，提供交互式 SSH 会话给 opencode 使用  
**目标用户**: 开发者/运维人员，需要通过 opencode 管理多台远程 Linux 服务器  
**技术栈**: Go (轻量级，编译成单一可执行文件)

## 特性

- **MCP 协议兼容**: 符合 Model Context Protocol 标准
- **SSH Config 支持**: 解析 `~/.ssh/config` 和 `~/.ssh/config.d/*.conf`
- **会话管理**: 保持连接状态，连续操作无需重复指定目标
- **多主机管理**: 轻松在多个服务器之间切换
- **安全认证**: 支持 SSH 密钥认证，无需密码
- **自动重试**: 连接失败时自动重试
- **轻量级**: Go 编译的单一二进制文件，资源占用少

## 安装

### 前置要求
- Go 1.21+
- opencode 或其他 MCP 兼容客户端
- SSH 配置文件 (`~/.ssh/config`) 已配置好目标服务器

### 安装方式
```bash
# 从源码安装
go install github.com/quotar/opencode-ssh-mcp@latest

# 或下载预编译二进制（发布版）
wget https://github.com/quotar/opencode-ssh-mcp/releases/latest/download/opencode-ssh-mcp_linux_amd64
chmod +x opencode-ssh-mcp_linux_amd64
```

## 配置

在 opencode 配置文件中（例如 `~/.opencode.json`）添加：

```json
{
  "mcpServers": {
    "ssh-manager": {
      "type": "local",
      "command": "/path/to/opencode-ssh-mcp",
      "args": []
    }
  }
}
```

## 使用

### opencode 中的使用
```
# 列出所有配置的 SSH 主机
> ssh list
Available hosts:
  - example-server (203.0.113.10)
  - test-server (203.0.113.11)

# 连接到指定主机
> ssh connect example-server
Connected to example-server

# 连接后直接执行命令（自动发送到远程）
> ls -la
[远程服务器输出]

> df -h
[远程服务器输出]

# 查看当前连接状态
> ssh status
Active connections: 1
Current active host: example-server

# 斷开连接
> ssh disconnect
Disconnected from example-server
```

## 架构

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

## 为什么选择 MCP 而不是 CLI

在 SSH 服务器管理场景中，MCP 比传统 CLI 方式更优：

- **Token 消耗**: 减少 70%+ (连接后直接执行命令 vs 每次 `ssh host 'command'`)
- **会话管理**: 保持连接状态，AI Agent 清楚当前操作目标
- **连续操作**: 一次连接，多次操作，避免重复连接开销
- **AI 体验**: 类似人类 SSH 会话的自然交互

## 贡献

参见 [CONTRIBUTING.md](CONTRIBUTING.md)

## 许可证

MIT - 参见 [LICENSE](LICENSE)