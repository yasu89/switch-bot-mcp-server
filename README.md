# SwitchBot MCP Server

[日本語はこちら](README_ja.md)

The SwitchBot MCP Server is a [Model Context Protocol (MCP)](https://modelcontextprotocol.io/introduction) server that provides a feature to control SwitchBot devices interactively using [SwitchBotAPI](https://github.com/OpenWonderLabs/SwitchBotAPI).

## Use Cases

- Operate SwitchBot devices interactively
- Perform operations on multiple devices at once
- Use data retrieved from a device to operate another device

## Installation

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

Currently, only a few tools are available, such as retrieving devices, retrieving statuses, and executing ON/OFF commands.

| Tool Name                      | Description                 |
|--------------------------------|-----------------------------|
| `get_switch_bot_devices`       | Get SwitchBot devices       |
| `get_switch_bot_device_status` | Get SwitchBot device status |
| `turn_on_off_device`           | Turn on/off device          |
