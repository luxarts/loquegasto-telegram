package controller

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParserController_GetTypeFromMessage_TypePayment(t *testing.T) {
	// Given
	msg := "123 some item"
	ctrl := NewParserController()

	// When
	msgType := ctrl.GetTypeFromMessage(msg)

	// Then
	require.Equal(t, MessageTypePayment, msgType)
}
func TestParserController_GetTypeFromMessage_Unknown(t *testing.T) {
	// Given
	msg := "abcdefg"
	ctrl := NewParserController()

	// When
	msgType := ctrl.GetTypeFromMessage(msg)

	// Then
	require.Equal(t, MessageTypeUnknown, msgType)
}
