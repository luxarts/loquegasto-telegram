package controller

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParserController_GetTypeFromMessage_TypePaymentSimple(t *testing.T) {
	// Given
	msg := "1234 some item"
	ctrl := NewParserController(nil, nil)

	// When
	msgType := ctrl.GetTypeFromMessage(msg)

	// Then
	require.EqualValues(t, messageTypePayment, msgType)
}
func TestParserController_GetTypeFromMessage_TypePaymentWithSource(t *testing.T) {
	// Given
	msg := "123 some item (source)"
	ctrl := NewParserController(nil, nil)

	// When
	msgType := ctrl.GetTypeFromMessage(msg)

	// Then
	require.EqualValues(t, messageTypePayment, msgType)
}
func TestParserController_GetTypeFromMessage_TypePaymentWithSign(t *testing.T) {
	// Given
	msg := "$123 some item"
	ctrl := NewParserController(nil, nil)

	// When
	msgType := ctrl.GetTypeFromMessage(msg)

	// Then
	require.EqualValues(t, messageTypePayment, msgType)
}
func TestParserController_GetTypeFromMessage_TypePaymentWithFloatValue(t *testing.T) {
	// Given
	msg := "123.456 some item"
	ctrl := NewParserController(nil, nil)

	// When
	msgType := ctrl.GetTypeFromMessage(msg)

	// Then
	require.EqualValues(t, messageTypePayment, msgType)
}
func TestParserController_GetTypeFromMessage_Unknown(t *testing.T) {
	// Given
	msg := "abcdefg"
	ctrl := NewParserController(nil, nil)

	// When
	msgType := ctrl.GetTypeFromMessage(msg)

	// Then
	require.EqualValues(t, messageTypeUnknown, msgType)
}
