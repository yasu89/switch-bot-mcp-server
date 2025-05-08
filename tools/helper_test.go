package tools_test

import (
	"encoding/json"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yasu89/switch-bot-api-go"
	"reflect"
	"testing"
)

// createMCPRequest is a helper function to create a MCP request with the given arguments.
func createMCPRequest(args map[string]interface{}) mcp.CallToolRequest {
	return mcp.CallToolRequest{
		Params: struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
			Meta      *struct {
				ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
			} `json:"_meta,omitempty"`
		}{
			Arguments: args,
		},
	}
}

// getTextResult is a helper function that returns a text result from a tool call.
func getTextResult(t *testing.T, result *mcp.CallToolResult) mcp.TextContent {
	t.Helper()
	assert.NotNil(t, result)
	require.Len(t, result.Content, 1)
	require.IsType(t, mcp.TextContent{}, result.Content[0])
	textContent := result.Content[0].(mcp.TextContent)
	assert.Equal(t, "text", textContent.Type)
	return textContent
}

// assertDevicesResponse asserts that the given GetDevicesResponse matches the expected list of devices.
func assertDevicesResponse(t *testing.T, response *switchbot.GetDevicesResponseBody, expectedDeviceList []interface{}, expectedInfraredList []interface{}) {
	t.Helper()

	if len(response.DeviceList) != len(expectedDeviceList) {
		t.Fatalf("Expected %d device, got %d", len(expectedDeviceList), len(response.DeviceList))
	}

	if len(response.InfraredRemoteList) != len(expectedInfraredList) {
		t.Fatalf("Expected %d infrared device, got %d", len(expectedInfraredList), len(response.InfraredRemoteList))
	}

	for i, device := range response.DeviceList {
		jsonText, err := json.Marshal(expectedDeviceList[i])
		if err != nil {
			t.Fatalf("Failed to marshal expected device: %v", err)
		}
		var expectedDeviceMap map[string]interface{}
		err = json.Unmarshal(jsonText, &expectedDeviceMap)
		if err != nil {
			t.Fatalf("Failed to unmarshal expected device: %v", err)
		}

		if !reflect.DeepEqual(device, expectedDeviceMap) {
			t.Fatalf("expected %s, actual %s", jsonDump(t, expectedDeviceMap), jsonDump(t, device))
		}
	}

	for i, device := range response.InfraredRemoteList {
		jsonText, err := json.Marshal(expectedInfraredList[i])
		if err != nil {
			t.Fatalf("Failed to marshal expected device: %v", err)
		}
		var expectedDeviceMap map[string]interface{}
		err = json.Unmarshal(jsonText, &expectedDeviceMap)
		if err != nil {
			t.Fatalf("Failed to unmarshal expected device: %v", err)
		}

		if !reflect.DeepEqual(device, expectedDeviceMap) {
			t.Fatalf("expected %s, actual %s", jsonDump(t, expectedDeviceMap), jsonDump(t, device))
		}
	}
}

// jsonDump is a helper function to pretty-print JSON data for debugging.
func jsonDump(t *testing.T, data interface{}) string {
	t.Helper()

	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	return string(b)
}
