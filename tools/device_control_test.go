package tools_test

import (
	"context"
	"github.com/yasu89/switch-bot-api-go/helpers"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yasu89/switch-bot-api-go"
	"github.com/yasu89/switch-bot-mcp-server/tools"
)

func Test_GetExecuteCommandTool(t *testing.T) {
	t.Run("Test tool", func(t *testing.T) {
		client := switchbot.NewClient("secret", "token")
		tool, _ := tools.GetExecuteCommandTool(client)

		assert.Equal(t, "execute_command", tool.Name)
		assert.NotEmpty(t, tool.Description)
		assert.Contains(t, tool.InputSchema.Properties, "device_id")
		assert.Contains(t, tool.InputSchema.Properties, "command_parameter_json")
		assert.ElementsMatch(t, tool.InputSchema.Required, []string{"device_id", "command_parameter_json"})
	})

	physicalBotDeviceID := "BOT123456"
	physicalCeilingLightDeviceID := "CEILINGLIGHT123456"
	infraredDeviceId := "AIRCONDITIONER654321"

	testDataList := []struct {
		name             string
		expectedBody     string
		expectedDeviceID string
		args             map[string]interface{}
	}{
		{
			name:             "Test handler physical BotDevice TurnOn",
			expectedBody:     `{"commandType": "command","command": "turnOn","parameter": "default"}`,
			expectedDeviceID: physicalBotDeviceID,
			args: map[string]interface{}{
				"device_id":              physicalBotDeviceID,
				"command_parameter_json": `{"command":"TurnOn"}`,
			},
		},
		{
			name:             "Test handler physical BotDevice TurnOff",
			expectedBody:     `{"commandType": "command","command": "turnOff","parameter": "default"}`,
			expectedDeviceID: physicalBotDeviceID,
			args: map[string]interface{}{
				"device_id":              physicalBotDeviceID,
				"command_parameter_json": `{"command":"TurnOff"}`,
			},
		},
		{
			name:             "Test handler physical BotDevice SetBrightness",
			expectedBody:     `{"commandType": "command","command": "setBrightness","parameter": "50"}`,
			expectedDeviceID: physicalCeilingLightDeviceID,
			args: map[string]interface{}{
				"device_id":              physicalCeilingLightDeviceID,
				"command_parameter_json": `{"command":"SetBrightness", "brightness": 50}`,
			},
		},
		{
			name:             "Test handler infrared device TurnOn",
			expectedBody:     `{"commandType": "command","command": "turnOn","parameter": "default"}`,
			expectedDeviceID: infraredDeviceId,
			args: map[string]interface{}{
				"device_id":              infraredDeviceId,
				"command_parameter_json": `{"command":"TurnOn"}`,
			},
		},
		{
			name:             "Test handler infrared device TurnOff",
			expectedBody:     `{"commandType": "command","command": "turnOff","parameter": "default"}`,
			expectedDeviceID: infraredDeviceId,
			args: map[string]interface{}{
				"device_id":              infraredDeviceId,
				"command_parameter_json": `{"command":"TurnOff"}`,
			},
		},
		{
			name:             "Test handler air conditioner device SetAll",
			expectedBody:     `{"commandType": "command","command": "setAll","parameter": "20,2,2,on"}`,
			expectedDeviceID: infraredDeviceId,
			args: map[string]interface{}{
				"device_id":              infraredDeviceId,
				"command_parameter_json": `{"command": "SetAll", "temperatureCelsius": 20, "mode": 2, "fan": 2, "powerState": "on"}`,
			},
		},
	}

	for _, testData := range testDataList {
		t.Run(testData.name, func(t *testing.T) {
			switchBotMock := helpers.NewSwitchBotMock(t)
			switchBotMock.RegisterDevicesMock(
				[]interface{}{
					map[string]interface{}{
						"deviceId":           physicalBotDeviceID,
						"deviceType":         "Bot",
						"hubDeviceId":        "123456789",
						"deviceName":         "BotDevice",
						"enableCloudService": true,
					},
					map[string]interface{}{
						"deviceId":           physicalCeilingLightDeviceID,
						"deviceType":         "Ceiling Light",
						"hubDeviceId":        "123456789",
						"deviceName":         "CeilingLightDevice",
						"enableCloudService": true,
					},
				},
				[]interface{}{
					map[string]interface{}{
						"deviceId":    infraredDeviceId,
						"deviceName":  "AirConditionerDevice",
						"remoteType":  "Air Conditioner",
						"hubDeviceId": "HUB123456",
					},
				},
			)
			switchBotMock.RegisterCommandMock(testData.expectedDeviceID, testData.expectedBody)
			testServer := switchBotMock.NewTestServer()
			defer testServer.Close()

			client := switchbot.NewClient("secret", "token", switchbot.OptionBaseApiURL(testServer.URL))
			_, handler := tools.GetExecuteCommandTool(client)

			request := createMCPRequest(testData.args)
			_, err := handler(context.Background(), request)
			assert.NoError(t, err)

			switchBotMock.AssertCallCount(http.MethodGet, "/devices", 1)
			switchBotMock.AssertCallCount(http.MethodPost, "/devices/"+testData.expectedDeviceID+"/commands", 1)
		})
	}
}
