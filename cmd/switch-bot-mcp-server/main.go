package main

import (
	"fmt"
	"github.com/mark3labs/mcp-go/server"
	"github.com/yasu89/switch-bot-mcp-server/tools"
	"github.com/yasu89/switch-bot-mcp-server/version"
	"log"
	"os"

	"github.com/yasu89/switch-bot-api-go"
)

func main() {
	token, ok := os.LookupEnv("SWITCH_BOT_TOKEN")
	if !ok {
		log.Fatal("SWITCH_BOT_TOKEN environment variable is required")
	}
	secret, ok := os.LookupEnv("SWITCH_BOT_SECRET")
	if !ok {
		log.Fatal("SWITCH_BOT_SECRET environment variable is required")
	}

	switchBotClient := switchbot.NewClient(secret, token)

	// Create a new MCP server
	mcpServer := server.NewMCPServer(
		"SwitchBot MCP",
		version.Version,
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	mcpServer.AddTool(tools.GetDeviceListTool(switchBotClient))
	mcpServer.AddTool(tools.GetDeviceStatusTool(switchBotClient))
	mcpServer.AddTool(tools.GetExecuteCommandTool(switchBotClient))

	// Start the server
	if err := server.ServeStdio(mcpServer); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
