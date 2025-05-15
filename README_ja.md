# SwitchBot MCP Server

[![Go Report Card](https://goreportcard.com/badge/github.com/yasu89/switch-bot-mcp-server)](https://goreportcard.com/report/github.com/yasu89/switch-bot-mcp-server)
![Coverage](https://raw.githubusercontent.com/yasu89/octocovs/main/badges/yasu89/switch-bot-mcp-server/coverage.svg)
![Code to Test Ratio](https://raw.githubusercontent.com/yasu89/octocovs/main/badges/yasu89/switch-bot-mcp-server/ratio.svg)
![Test Execution Time](https://raw.githubusercontent.com/yasu89/octocovs/main/badges/yasu89/switch-bot-mcp-server/time.svg)

[English version is here.](README.md)

SwitchBot MCP Serverは[SwitchBotAPI](https://github.com/OpenWonderLabs/SwitchBotAPI)を使用してSwitchBotのデバイスを会話で操作できる機能を提供する[Model Context Protocol (MCP)](https://modelcontextprotocol.io/introduction)をサーバです。

## 使用例

- SwitchBotのデバイスを対話を通して操作する
- 複数のデバイスに対する操作を一度に行う
- デバイスから取得したデータを元に別のデバイスの操作を行う

## インストール方法

### バイナリ準備

[リリースページ](https://github.com/yasu89/switch-bot-mcp-server/releases)からダウンロードしてください。

### シークレットとトークンの準備

[SwitchBotAPIのGetting Started](https://github.com/OpenWonderLabs/SwitchBotAPI?tab=readme-ov-file#getting-started)に従って、SwitchBotAPIのトークンとシークレットを取得してください。

### Claude Desktopで使用する場合の設定

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

## 利用可能なツール

デバイスの取得とステータスの取得、デバイスのコマンドの実行が利用可能です。

| Tool Name                      | Description                   |
|--------------------------------|-------------------------------|
| `get_switch_bot_devices`       | Get SwitchBot devices         |
| `get_switch_bot_device_status` | Get SwitchBot device status   |
| `execute_command`              | Execute a command on a device |
