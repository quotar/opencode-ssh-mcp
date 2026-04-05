# 常见问题解答

## 通用问题

### Q: 这个项目是做什么的?

A: opencode-ssh-mcp 是一个 MCP (Model Context Protocol) 服务器，它允许 AI Agent (如 opencode) 通过 SSH 协议管理远程服务器，提供"本地操作远程"的体验。

### Q: 为什么需要这个工具而不是直接使用 SSH?

A: 直接使用 SSH 有以下挑战：
- 每个命令都需要完整的 `ssh user@host 'command'` 格式
- 对 AI Agent 来说心智负担重
- token 消耗高
- 无状态管理

MCP 方式提供会话管理，连接后可直接执行命令。

### Q: 如何配置 SSH 连接?

A: 确保 `~/.ssh/config` 文件配置好目标服务器：

```
Host myserver
    HostName 203.0.113.10
    User myuser
    IdentityFile ~/.ssh/id_rsa
```

## MCP 协议相关

### Q: 什么是 MCP?

A: MCP (Model Context Protocol) 是为 AI Agent 提供外部系统访问的标准协议，类似于为 AI 提供 USB-C 接口。

### Q: 为什么使用 MCP 而不是 CLI?

A: 在 SSH 管理场景中，MCP 提供：
- 会话状态管理
- 减少 70%+ token 消耗
- 连接复用
- 自然的交互体验

## 安装和配置

### Q: 如何安装?

A: 从源码安装：
```bash
go install github.com/quotar/opencode-ssh-mcp@latest
```

或下载预编译二进制文件。

### Q: 如何配置 opencode?

A: 在 opencode 配置文件中添加：

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

## 使用问题

### Q: 如何使用?

A: 在 opencode 中：
1. `ssh list` - 列出可用主机
2. `ssh connect <host>` - 连接到主机
3. 直接执行命令（自动发送到远程）
4. `ssh disconnect` - 斷开

### Q: 支持多少个服务器?

A: 支持任意数量的服务器，只要它们在 SSH 配置文件中定义。

### Q: 如何切换服务器?

A: 使用 `ssh connect <new_host>` 切换到不同服务器。

## 安全问题

### Q: 这个工具安全吗?

A: 是的，使用以下安全措施：
- 标准 SSH 库 (golang.org/x/crypto/ssh)
- 支持密钥认证
- 输入验证
- 无额外权限要求

### Q: 需要存储密码吗?

A: 不需要，支持 SSH 密钥认证，无需存储密码。

### Q: 如何处理认证失败?

A: 实现自动重试机制，失败时返回详细错误信息。

## 故障排除

### Q: 连接失败怎么办?

A: 检查：
1. SSH 配置文件 (`~/.ssh/config`)
2. 密钥权限 (600 for keys)
3. 目标服务器可达性
4. 防火墙设置

### Q: 命令不执行怎么办?

A: 确认：
1. 已成功连接到主机 (`ssh status`)
2. 命令格式正确
3. 目标服务器有相应权限

### Q: 如何查看当前状态?

A: 使用 `ssh status` 命令查看当前连接状态。

## 性能问题

### Q: 会影响性能吗?

A: 不会，Go 编译的单一二进制文件，资源占用少 (~20MB 内存)，启动快。

### Q: 有连接池吗?

A: 是的，支持连接复用和管理。

## 开发问题

### Q: 如何为这项目贡献?

A: 参见 CONTRIBUTING.md 文件中的贡献指南。

### Q: 支持哪些平台?

A: Linux, macOS, Windows (通过预编译二进制文件)。