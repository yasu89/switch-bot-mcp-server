# SwitchBot MCP Server

[![Go Report Card](https://goreportcard.com/badge/github.com/yasu89/switch-bot-mcp-server)](https://goreportcard.com/report/github.com/yasu89/switch-bot-mcp-server)
![Coverage](https://raw.githubusercontent.com/yasu89/octocovs/main/badges/yasu89/switch-bot-mcp-server/coverage.svg)
![Code to Test Ratio](https://raw.githubusercontent.com/yasu89/octocovs/main/badges/yasu89/switch-bot-mcp-server/ratio.svg)
![Test Execution Time](https://raw.githubusercontent.com/yasu89/octocovs/main/badges/yasu89/switch-bot-mcp-server/time.svg)

[日本語はこちら](README_ja.md)

The SwitchBot MCP Server is a [Model Context Protocol (MCP)](https://modelcontextprotocol.io/introduction) server that provides a feature to control SwitchBot devices interactively using [SwitchBotAPI](https://github.com/OpenWonderLabs/SwitchBotAPI).

## Use Cases

- Operate SwitchBot devices interactively
- Perform operations on multiple devices at once
- Use data retrieved from a device to operate another device

## Installation

### Prepare binary

Download binary from [release page](https://github.com/yasu89/switch-bot-mcp-server/releases).

### Prepare secret and token

Follow the [Getting Started guide of SwitchBotAPI](https://github.com/OpenWonderLabs/SwitchBotAPI?tab=readme-ov-file#getting-started) to obtain the token and secret for SwitchBotAPI.

### Setting for Claude Desktop

```json
{
  "mcpServers": {
    "switchbot": {
      "command": "~/Downloads/switch-bot-mcp-server",
      "env": {
        "SWITCH_BOT_TOKEN": "YOUR_SWITCH_BOT_TOKEN",
        "SWITCH_BOT_SECRET": "YOUR_SWITCH_BOT_SECRET"
      }
    }
  }
}
```

## Available Tools

Retrieving devices, retrieving statuses, and executing commands on devices are available.

| Tool Name                      | Description                   |
|--------------------------------|-------------------------------|
| `get_switch_bot_devices`       | Get SwitchBot devices         |
| `get_switch_bot_device_status` | Get SwitchBot device status   |
| `execute_command`              | Execute a command on a device |
