package tools_test

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/yasu89/switch-bot-api-go"
	"github.com/yasu89/switch-bot-api-go/helpers"
	"github.com/yasu89/switch-bot-mcp-server/tools"
	"net/http"
	"testing"
)

func Test_GetDeviceListTool(t *testing.T) {
	t.Run("Test tool", func(t *testing.T) {
		client := switchbot.NewClient("secret", "token")
		tool, _ := tools.GetDeviceListTool(client)

		assert.Equal(t, "get_switch_bot_devices", tool.Name)
		assert.NotEmpty(t, tool.Description)
		assert.Empty(t, tool.InputSchema.Properties)
		assert.Empty(t, tool.InputSchema.Required)
	})

	t.Run("Test handler", func(t *testing.T) {
		switchBotMock := helpers.NewSwitchBotMock(t)
		switchBotMock.RegisterDevicesMock(
			[]interface{}{
				map[string]interface{}{
					"deviceId":           "PHYSICAL123456",
					"deviceType":         "Bot",
					"hubDeviceId":        "123456789",
					"deviceName":         "BotDevice",
					"enableCloudService": true,
				},
			},
			[]interface{}{
				map[string]interface{}{
					"deviceId":    "INFRARED123456",
					"deviceName":  "AirConditionerDevice",
					"remoteType":  "Air Conditioner",
					"hubDeviceId": "HUB123456",
				},
			},
		)
		testServer := switchBotMock.NewTestServer()
		defer testServer.Close()

		client := switchbot.NewClient("secret", "token", switchbot.OptionBaseApiURL(testServer.URL))
		_, handler := tools.GetDeviceListTool(client)

		request := createMCPRequest(map[string]interface{}{})
		result, err := handler(context.Background(), request)
		assert.NoError(t, err)

		textContent := getTextResult(t, result)

		var devicesResponse switchbot.GetDevicesResponseBody
		err = json.Unmarshal([]byte(textContent.Text), &devicesResponse)
		assert.NoError(t, err)

		assertDevicesResponse(
			t, &devicesResponse,
			[]interface{}{
				&switchbot.BotDevice{
					CommonDeviceListItem: switchbot.CommonDeviceListItem{
						CommonDevice: switchbot.CommonDevice{
							DeviceID:    "PHYSICAL123456",
							DeviceType:  "Bot",
							HubDeviceId: "123456789",
						},
						Client:             client,
						DeviceName:         "BotDevice",
						EnableCloudService: true,
					},
				},
			},
			[]interface{}{
				&switchbot.InfraredRemoteAirConditionerDevice{
					InfraredRemoteDevice: switchbot.InfraredRemoteDevice{
						Client:      client,
						DeviceID:    "INFRARED123456",
						DeviceName:  "AirConditionerDevice",
						RemoteType:  "Air Conditioner",
						HubDeviceId: "HUB123456",
					},
				},
			},
		)

		switchBotMock.AssertCallCount(http.MethodGet, "/devices", 1)
	})
}
