package tools

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/yasu89/switch-bot-api-go"
)

// GetExecuteCommandTool creates a tool to execute a command on a specific SwitchBot device.
func GetExecuteCommandTool(switchBotClient *switchbot.Client) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"execute_command",
			mcp.WithDescription("Execute a command on a device"),
			mcp.WithString(
				"device_id",
				mcp.Required(),
				mcp.Description("ID of the device to execute a command on"),
			),
			mcp.WithString(
				"command_parameter_json",
				mcp.Required(),
				mcp.Description("Command to send"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			deviceId := request.Params.Arguments["device_id"].(string)
			commandParameterJsonText := request.Params.Arguments["command_parameter_json"].(string)

			response, err := switchBotClient.GetDevices()
			if err != nil {
				return nil, err
			}

			var targetDevice switchbot.ExecutableCommandDevice

			for _, physicalDevice := range response.Body.DeviceList {
				switch physicalDevice.(type) {
				case switchbot.ExecutableCommandDevice:
					deviceIDGettable := physicalDevice.(switchbot.DeviceIDGettable)
					executableCommandDevice := physicalDevice.(switchbot.ExecutableCommandDevice)
					if deviceIDGettable.GetDeviceID() != deviceId {
						continue
					}
					targetDevice = executableCommandDevice
				}
			}

			for _, infraredDevice := range response.Body.InfraredRemoteList {
				switch infraredDevice.(type) {
				case switchbot.ExecutableCommandDevice:
					deviceIDGettable := infraredDevice.(switchbot.DeviceIDGettable)
					executableCommandDevice := infraredDevice.(switchbot.ExecutableCommandDevice)
					if deviceIDGettable.GetDeviceID() != deviceId {
						continue
					}
					targetDevice = executableCommandDevice
				}
			}

			if targetDevice == nil {
				return mcp.NewToolResultError("Device not found"), nil
			}

			commandResponse, err := targetDevice.ExecCommand(commandParameterJsonText)
			if err != nil {
				return nil, err
			}

			responseJsonText, err := json.Marshal(commandResponse)
			if err != nil {
				return nil, err
			}

			return mcp.NewToolResultText(string(responseJsonText)), nil
		}
}
