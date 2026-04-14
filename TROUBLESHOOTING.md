# opencode-ssh-mcp 集成问题记录

## 问题概述

在尝试将 opencode-ssh-mcp 集成到 opencode 1.3.17 中时，发现设计不匹配导致无法正常使用MCP协议进行远程机器操作。

## 详细问题描述

### 1. 配置兼容性问题
- opencode 1.3.17 拒绝识别配置文件中的 `mcpServers` 键
- 添加该键会导致启动错误：`Unrecognized key: "mcpServers"`
- 这与 opencode-ssh-mcp 项目文档中的说明不符

### 2. MCP 服务器注册机制缺失
- `opencode mcp add` 命令不接受任何参数
- 无法通过命令行添加 SSH-MCP 服务器
- 未找到有效的配置入口来注册 MCP 服务器

### 3. 架构不匹配
- opencode-ssh-mcp 的设计与当前 opencode 版本的 MCP 机制不兼容
- 需要更新 opencode-ssh-mcp 以符合 opencode 的 MCP 注册协议

## 当前状态

- **SSH-MCP 项目功能**：✅ 正常工作，能解析 SSH 配置并连接服务器
- **opencode 集成**：❌ 失败，无法将 SSH-MCP 注册为 MCP 服务器
- **直接 SSH 连接**：✅ 工作正常，可通过 `ssh image-91` 直接连接

## 测试结果

### SSH-MCP 服务器测试
```bash
echo '{"jsonrpc":"2.0","id":"test","method":"initialize","params":{}}' | /opt/github/opencode-ssh-mcp/opencode-ssh-mcp
```
结果：✅ 成功，加载了 40 个 SSH 配置包括 image-91

### opencode 配置测试
添加 `mcpServers` 后：
```bash
opencode
```
结果：❌ 失败，`Unrecognized key: "mcpServers"`

## 解决方案建议

1. 修改 opencode-ssh-mcp 以符合 opencode 的 MCP 服务器注册协议
2. 更新 opencode 配置处理逻辑以支持 MCP 服务器注册
3. 或者创建兼容的插件机制来桥接两者

## 验证命令

```bash
# 验证 SSH-MCP 服务器功能
echo '{"jsonrpc":"2.0","id":"test","method":"initialize","params":{}}' | /opt/github/opencode-ssh-mcp/opencode-ssh-mcp

# 验证 opencode 配置兼容性
opencode debug config

# 验证直接 SSH 连接
ssh image-91 'hostname && whoami'
```

## 影响

目前无法通过 opencode AI Agent 自动化 SSH 操作，只能通过传统 SSH 命令手动连接服务器。AI Agent 无法使用 MCP 协议对远程服务器进行操作。