package tools

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/yasu89/switch-bot-api-go"
)

// AddTurnOnOffDeviceTool adds a tool to the MCP server that turns on/off a specific SwitchBot device.
func AddTurnOnOffDeviceTool(mcpServer *server.MCPServer, switchBotClient *switchbot.Client) {
	mcpServer.AddTool(
		mcp.NewTool(
			"turn_on_off_device",
			mcp.WithDescription("Turn on/off device"),
			mcp.WithString(
				"device_id",
				mcp.Required(),
				mcp.Description("ID of the device to turn on/off"),
			),
			mcp.WithBoolean(
				"is_turn_on",
				mcp.Required(),
				mcp.Description("Command to send (true:on, false:off)"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			deviceId := request.Params.Arguments["device_id"].(string)
			isTurnOn := request.Params.Arguments["is_turn_on"].(bool)

			response, err := switchBotClient.GetDevices()
			if err != nil {
				return nil, err
			}

			var targetDevice switchbot.SwitchableDevice

			for _, physicalDevice := range response.Body.DeviceList {
				switch physicalDevice.(type) {
				case switchbot.SwitchableDevice:
					deviceIDGettable := physicalDevice.(switchbot.DeviceIDGettable)
					switchableDevice := physicalDevice.(switchbot.SwitchableDevice)
					if deviceIDGettable.GetDeviceID() != deviceId {
						continue
					}
					targetDevice = switchableDevice
				}
			}

			for _, infraredDevice := range response.Body.InfraredRemoteList {
				switch infraredDevice.(type) {
				case switchbot.SwitchableDevice:
					deviceIDGettable := infraredDevice.(switchbot.DeviceIDGettable)
					switchableDevice := infraredDevice.(switchbot.SwitchableDevice)
					if deviceIDGettable.GetDeviceID() != deviceId {
						continue
					}
					targetDevice = switchableDevice
				}
			}

			if targetDevice != nil {
				var commandResponse *switchbot.CommonResponse
				if isTurnOn {
					commandResponse, err = targetDevice.TurnOn()
				} else {
					commandResponse, err = targetDevice.TurnOff()
				}

				if err != nil {
					return nil, err
				}

				responseJsonText, err := json.Marshal(commandResponse)
				if err != nil {
					return nil, err
				}

				return mcp.NewToolResultText(string(responseJsonText)), nil
			}

			return mcp.NewToolResultError("Device not found"), nil
		},
	)
}
