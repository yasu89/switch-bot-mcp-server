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

func Test_GetTurnOnOffDeviceTool(t *testing.T) {
	t.Run("Test tool", func(t *testing.T) {
		client := switchbot.NewClient("secret", "token")
		tool, _ := tools.GetTurnOnOffDeviceTool(client)

		assert.Equal(t, "turn_on_off_device", tool.Name)
		assert.NotEmpty(t, tool.Description)
		assert.Contains(t, tool.InputSchema.Properties, "device_id")
		assert.Contains(t, tool.InputSchema.Properties, "is_turn_on")
		assert.ElementsMatch(t, tool.InputSchema.Required, []string{"device_id", "is_turn_on"})
	})

	physicalDeviceID := "ABCDEF123456"
	infraredDeviceId := "ABCDEF654321"

	testDataList := []struct {
		name             string
		expectedBody     string
		expectedDeviceID string
		args             map[string]interface{}
	}{
		{
			name:             "Test handler physical device TurnOn",
			expectedBody:     "{\"commandType\": \"command\",\"command\": \"turnOn\",\"parameter\": \"default\"}",
			expectedDeviceID: physicalDeviceID,
			args: map[string]interface{}{
				"device_id":  physicalDeviceID,
				"is_turn_on": true,
			},
		},
		{
			name:             "Test handler physical device TurnOff",
			expectedBody:     "{\"commandType\": \"command\",\"command\": \"turnOff\",\"parameter\": \"default\"}",
			expectedDeviceID: physicalDeviceID,
			args: map[string]interface{}{
				"device_id":  physicalDeviceID,
				"is_turn_on": false,
			},
		},
		{
			name:             "Test handler infrared device TurnOn",
			expectedBody:     "{\"commandType\": \"command\",\"command\": \"turnOn\",\"parameter\": \"default\"}",
			expectedDeviceID: infraredDeviceId,
			args: map[string]interface{}{
				"device_id":  infraredDeviceId,
				"is_turn_on": true,
			},
		},
		{
			name:             "Test handler infrared device TurnOff",
			expectedBody:     "{\"commandType\": \"command\",\"command\": \"turnOff\",\"parameter\": \"default\"}",
			expectedDeviceID: infraredDeviceId,
			args: map[string]interface{}{
				"device_id":  infraredDeviceId,
				"is_turn_on": false,
			},
		},
	}

	for _, testData := range testDataList {
		t.Run(testData.name, func(t *testing.T) {
			switchBotMock := helpers.NewSwitchBotMock(t)
			switchBotMock.RegisterDevicesMock(
				[]interface{}{
					map[string]interface{}{
						"deviceId":           physicalDeviceID,
						"deviceType":         "Bot",
						"hubDeviceId":        "123456789",
						"deviceName":         "BotDevice",
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
			_, handler := tools.GetTurnOnOffDeviceTool(client)

			request := createMCPRequest(testData.args)

			_, err := handler(context.Background(), request)
			assert.NoError(t, err)

			switchBotMock.AssertCallCount(http.MethodPost, "/devices/"+testData.expectedDeviceID+"/commands", 1)
		})
	}
}
