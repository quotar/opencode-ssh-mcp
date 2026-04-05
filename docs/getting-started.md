# 快速开始

## 安装

### 前置要求
- Go 1.21+
- opencode 或其他 MCP 兼容客户端
- SSH 配置 (`~/.ssh/config`) 已设置

### 从源码安装
```bash
git clone https://github.com/quotar/opencode-ssh-mcp.git
cd opencode-ssh-mcp
make build
```

### 二进制安装
下载对应平台的预编译二进制文件。

## 配置

### SSH 配置
确保 `~/.ssh/config` 文件包含目标服务器：

```
Host example-server
    HostName 203.0.113.10
    User myuser
    IdentityFile ~/.ssh/id_rsa
```

### opencode 配置
在 opencode 配置文件中添加：

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

### 连接服务器
```
> ssh list
# 列出所有配置的服务器

> ssh connect example-server
# 连接到指定服务器
```

### 执行命令
```
# 连接后，直接执行命令
> ls -la
> df -h
> ps aux
```

### 状态管理
```
> ssh status
# 查看当前连接状态

> ssh disconnect
# 斷开连接
```