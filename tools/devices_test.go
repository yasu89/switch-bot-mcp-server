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

func Test_GetDeviceStatusTool(t *testing.T) {
	t.Run("Test tool", func(t *testing.T) {
		client := switchbot.NewClient("secret", "token")
		tool, _ := tools.GetDeviceStatusTool(client)

		assert.Equal(t, "get_switch_bot_device_status", tool.Name)
		assert.NotEmpty(t, tool.Description)
		assert.Contains(t, tool.InputSchema.Properties, "device_id")
		assert.ElementsMatch(t, tool.InputSchema.Required, []string{"device_id"})
	})

	t.Run("Test handler", func(t *testing.T) {
		switchBotMock := helpers.NewSwitchBotMock(t)
		switchBotMock.RegisterDevicesMock(
			[]interface{}{
				map[string]interface{}{
					"deviceId":           "ABCDEF123456",
					"deviceType":         "Bot",
					"hubDeviceId":        "123456789",
					"deviceName":         "BotDevice",
					"enableCloudService": true,
				},
			},
			[]interface{}{},
		)
		switchBotMock.RegisterStatusMock("ABCDEF123456", map[string]interface{}{
			"deviceId":    "ABCDEF123456",
			"deviceType":  "Bot",
			"hubDeviceId": "123456789",
			"power":       "ON",
			"battery":     100,
			"version":     "1.0",
			"deviceMode":  "pressMode",
		})
		testServer := switchBotMock.NewTestServer()
		defer testServer.Close()

		client := switchbot.NewClient("secret", "token", switchbot.OptionBaseApiURL(testServer.URL))
		_, handler := tools.GetDeviceStatusTool(client)

		request := createMCPRequest(map[string]interface{}{"device_id": "ABCDEF123456"})
		result, err := handler(context.Background(), request)
		assert.NoError(t, err)

		textContent := getTextResult(t, result)

		var deviceStatusResponse switchbot.BotDeviceStatusBody
		err = json.Unmarshal([]byte(textContent.Text), &deviceStatusResponse)
		assert.NoError(t, err)

		expectedBody := &switchbot.BotDeviceStatusBody{
			CommonDevice: switchbot.CommonDevice{
				DeviceID:    "ABCDEF123456",
				DeviceType:  "Bot",
				HubDeviceId: "123456789",
			},
			Power:      "ON",
			Battery:    100,
			Version:    "1.0",
			DeviceMode: "pressMode",
		}

		assertBody(t, &deviceStatusResponse, expectedBody)

		switchBotMock.AssertCallCount(http.MethodGet, "/devices", 1)
		switchBotMock.AssertCallCount(http.MethodGet, "/devices/ABCDEF123456/status", 1)
	})
}
