package tools

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/yasu89/switch-bot-api-go"
)

// AddGetDevicesTool adds a tool to the MCP server that retrieves SwitchBot devices.
func AddGetDevicesTool(mcpServer *server.MCPServer, switchBotClient *switchbot.Client) {
	mcpServer.AddTool(
		mcp.NewTool(
			"get_switch_bot_devices",
			mcp.WithDescription("Get SwitchBot devices"),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			response, err := switchBotClient.GetDevices()
			if err != nil {
				return nil, err
			}

			responseJsonText, err := json.Marshal(response.Body)
			if err != nil {
				return nil, err
			}

			return mcp.NewToolResultText(string(responseJsonText)), nil
		},
	)
}

// AddGetDeviceStatusTool adds a tool to the MCP server that retrieves the status of a specific SwitchBot device.
func AddGetDeviceStatusTool(mcpServer *server.MCPServer, switchBotClient *switchbot.Client) {
	mcpServer.AddTool(
		mcp.NewTool(
			"get_switch_bot_device_status",
			mcp.WithDescription("Get SwitchBot device status"),
			mcp.WithString(
				"device_id",
				mcp.Required(),
				mcp.Description("ID of the device to retrieve the status for"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			deviceId := request.Params.Arguments["device_id"].(string)

			response, err := switchBotClient.GetDevices()
			if err != nil {
				return nil, err
			}

			for _, device := range response.Body.DeviceList {
				switch device.(type) {
				case switchbot.StatusGettable:
					deviceIDGettable := device.(switchbot.DeviceIDGettable)
					device := device.(switchbot.StatusGettable)
					if deviceIDGettable.GetDeviceID() != deviceId {
						continue
					}

					statusResponseBody, err := device.GetAnyStatusBody()
					if err != nil {
						return nil, err
					}
					responseJonText, err := json.Marshal(statusResponseBody)
					if err != nil {
						return nil, err
					}
					return mcp.NewToolResultText(string(responseJonText)), nil
				}
			}

			return mcp.NewToolResultError("Device not found"), nil
		},
	)
}
