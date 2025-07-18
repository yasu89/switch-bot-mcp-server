package tools

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/yasu89/switch-bot-api-go"
)

// GetDeviceListTool creates a tool to get the list of SwitchBot devices.
func GetDeviceListTool(switchBotClient *switchbot.Client) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"get_switch_bot_devices",
			mcp.WithDescription("Get SwitchBot devices"),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			response, err := switchBotClient.GetDevices()
			if err != nil {
				return nil, err
			}

			var devicesResponse switchbot.GetDevicesResponseBody

			for _, device := range response.Body.DeviceList {
				var deviceMap map[string]interface{}

				deviceJsonText, err := json.Marshal(device)
				if err != nil {
					return nil, err
				}
				err = json.Unmarshal(deviceJsonText, &deviceMap)
				if err != nil {
					return nil, err
				}

				switch device.(type) {
				case switchbot.ExecutableCommandDevice:
					jsonSchema, err := device.(switchbot.ExecutableCommandDevice).GetCommandParameterJSONSchema()
					if err != nil {
						return nil, err
					}
					deviceMap["commandParameterJSONSchema"] = jsonSchema
				}

				devicesResponse.DeviceList = append(devicesResponse.DeviceList, deviceMap)
			}

			for _, device := range response.Body.InfraredRemoteList {
				var deviceMap map[string]interface{}

				deviceJsonText, err := json.Marshal(device)
				if err != nil {
					return nil, err
				}
				err = json.Unmarshal(deviceJsonText, &deviceMap)
				if err != nil {
					return nil, err
				}

				switch device.(type) {
				case switchbot.ExecutableCommandDevice:
					jsonSchema, err := device.(switchbot.ExecutableCommandDevice).GetCommandParameterJSONSchema()
					if err != nil {
						return nil, err
					}
					deviceMap["commandParameterJSONSchema"] = jsonSchema
				}

				devicesResponse.InfraredRemoteList = append(devicesResponse.InfraredRemoteList, deviceMap)
			}

			responseJsonText, err := json.Marshal(devicesResponse)
			if err != nil {
				return nil, err
			}

			return mcp.NewToolResultText(string(responseJsonText)), nil
		}
}

// GetDeviceStatusTool create a tool to get the status of a specific SwitchBot device.
func GetDeviceStatusTool(switchBotClient *switchbot.Client) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
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
		}
}
