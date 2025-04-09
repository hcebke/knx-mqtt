package utils

import (
	knxgo "github.com/vapourismo/knx-go/knx"
)

// KNXCommandToString converts a KNX command to its string representation
func KNXCommandToString(command knxgo.GroupCommand) string {
	switch command {
	case knxgo.GroupRead:
		return "GroupValue_Read"
	case knxgo.GroupWrite:
		return "GroupValue_Write"
	case knxgo.GroupResponse:
		return "GroupValue_Response"
	default:
		return "unknown"
	}
}
