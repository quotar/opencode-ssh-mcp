# 架构设计

## 概述

opencode-ssh-mcp 是一个 MCP (Model Context Protocol) 服务器，允许 AI Agent 通过标准 MCP 协议管理远程 SSH 服务器。

## 系统架构

### 组件

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

## MCP 协议实现

### 工具定义
- `ssh_list` - 列出所有配置的 SSH 主机
- `ssh_connect` - 连接到指定主机
- `ssh_exec` - 在远程执行命令
- `ssh_status` - 查看连接状态
- `ssh_disconnect` - 斷开连接

### 协议交互
使用标准 JSON-RPC 2.0 协议，STDIO 传输。

## 内部组件

### SSH 配置解析器 (internal/ssh/config.go)
- 解析 `~/.ssh/config` 和 `~/.ssh/config.d/*.conf`
- 支持 Host, HostName, Port, User, IdentityFile, ProxyJump

### 连接管理器 (internal/ssh/connection.go)
- 管理 SSH 连接池
- 实现自动重试机制
- 处理认证和会话

### 会话管理器 (internal/session/manager.go)
- 管理活动会话状态
- 跟踪当前活动主机
- 处理多主机切换

### MCP 服务器 (internal/mcp/server.go)
- 处理 MCP 协议请求
- 管理工具调用
- 维护会话状态

## 安全考虑

- 使用标准的 `golang.org/x/crypto/ssh` 库
- 支持现代加密算法
- 输入验证和命令注入防护
- 会话超时和资源限制

## 性能优化

- 连接复用减少重连开销
- 会话状态缓存
- 并发连接管理
- 轻量级 Go 二进制文件