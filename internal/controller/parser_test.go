package controller

import (
	tg "gopkg.in/telebot.v3"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParserController_GetTypeFromMessage_TypeNoDecimalNoSignNoWallet(t *testing.T) {
	// Given
	msg := "123 No Decimal No Sign No Category No Wallet"
	ctrl := NewParserController(nil, nil, nil, nil)

	// When
	msgType := ctrl.GetTypeFromMessage(&tg.Message{Payload: msg})

	// Then
	require.EqualValues(t, messageTypeTransaction, msgType)
}
func TestParserController_GetTypeFromMessage_TypeDotDecimalNoSignNoWallet(t *testing.T) {
	// Given
	msg := "123.45 Dot Decimal No Sign No Category No Wallet"
	ctrl := NewParserController(nil, nil, nil, nil)

	// When
	msgType := ctrl.GetTypeFromMessage(&tg.Message{Payload: msg})

	// Then
	require.EqualValues(t, messageTypeTransaction, msgType)
}
func TestParserController_GetTypeFromMessage_TypeCommaDecimalNoSignNoWallet(t *testing.T) {
	// Given
	msg := "123,45 Comma Decimal No Sign No Category No Wallet"
	ctrl := NewParserController(nil, nil, nil, nil)

	// When
	msgType := ctrl.GetTypeFromMessage(&tg.Message{Payload: msg})

	// Then
	require.EqualValues(t, messageTypeTransaction, msgType)
}
func TestParserController_GetTypeFromMessage_TypeNoDecimalSignNoWallet(t *testing.T) {
	// Given
	msg := "$123 No Decimal No Sign No Category No Wallet"
	ctrl := NewParserController(nil, nil, nil, nil)

	// When
	msgType := ctrl.GetTypeFromMessage(&tg.Message{Payload: msg})

	// Then
	require.EqualValues(t, messageTypeTransaction, msgType)
}
func TestParserController_GetTypeFromMessage_TypeDotDecimalSignNoWallet(t *testing.T) {
	// Given
	msg := "$123.45 Dot Decimal No Sign No Category No Wallet"
	ctrl := NewParserController(nil, nil, nil, nil)

	// When
	msgType := ctrl.GetTypeFromMessage(&tg.Message{Payload: msg})

	// Then
	require.EqualValues(t, messageTypeTransaction, msgType)
}
func TestParserController_GetTypeFromMessage_TypeCommaDecimalSignNoWallet(t *testing.T) {
	// Given
	msg := "$123,45 Comma Decimal No Sign No Category No Wallet"
	ctrl := NewParserController(nil, nil, nil, nil)

	// When
	msgType := ctrl.GetTypeFromMessage(&tg.Message{Payload: msg})

	// Then
	require.EqualValues(t, messageTypeTransaction, msgType)
}
func TestParserController_GetTypeFromMessage_TypeNoDecimalNoSignWallet(t *testing.T) {
	// Given
	msg := "123 No Decimal No Sign No Category (Wallet)"
	ctrl := NewParserController(nil, nil, nil, nil)

	// When
	msgType := ctrl.GetTypeFromMessage(&tg.Message{Payload: msg})

	// Then
	require.EqualValues(t, messageTypeTransaction, msgType)
}
func TestParserController_GetTypeFromMessage_TypeDotDecimalNoSignWallet(t *testing.T) {
	// Given
	msg := "123.45 Dot Decimal No Sign No Category (Wallet)"
	ctrl := NewParserController(nil, nil, nil, nil)

	// When
	msgType := ctrl.GetTypeFromMessage(&tg.Message{Payload: msg})

	// Then
	require.EqualValues(t, messageTypeTransaction, msgType)
}
func TestParserController_GetTypeFromMessage_TypeCommaDecimalNoSignWallet(t *testing.T) {
	// Given
	msg := "123,45 Comma Decimal No Sign No Category (Wallet)"
	ctrl := NewParserController(nil, nil, nil, nil)

	// When
	msgType := ctrl.GetTypeFromMessage(&tg.Message{Payload: msg})

	// Then
	require.EqualValues(t, messageTypeTransaction, msgType)
}
func TestParserController_GetTypeFromMessage_TypeNoDecimalSignWallet(t *testing.T) {
	// Given
	msg := "$123 No Decimal No Sign No Category (Wallet)"
	ctrl := NewParserController(nil, nil, nil, nil)

	// When
	msgType := ctrl.GetTypeFromMessage(&tg.Message{Payload: msg})

	// Then
	require.EqualValues(t, messageTypeTransaction, msgType)
}
func TestParserController_GetTypeFromMessage_TypeDotDecimalSignWallet(t *testing.T) {
	// Given
	msg := "$123.45 Dot Decimal No Sign No Category (Wallet)"
	ctrl := NewParserController(nil, nil, nil, nil)

	// When
	msgType := ctrl.GetTypeFromMessage(&tg.Message{Payload: msg})

	// Then
	require.EqualValues(t, messageTypeTransaction, msgType)
}
func TestParserController_GetTypeFromMessage_TypeCommaDecimalSignWallet(t *testing.T) {
	// Given
	msg := "$123,45 Comma Decimal No Sign No Category (Wallet)"
	ctrl := NewParserController(nil, nil, nil, nil)

	// When
	msgType := ctrl.GetTypeFromMessage(&tg.Message{Payload: msg})

	// Then
	require.EqualValues(t, messageTypeTransaction, msgType)
}
func TestParserController_GetTypeFromMessage_TypeNegativeNoDecimalNoSignNoWallet(t *testing.T) {
	// Given
	msg := "-123 Negative No Decimal No Sign No Wallet No Category"
	ctrl := NewParserController(nil, nil, nil, nil)

	// When
	msgType := ctrl.GetTypeFromMessage(&tg.Message{Payload: msg})

	// Then
	require.EqualValues(t, messageTypeTransaction, msgType)
}
func TestParserController_GetTypeFromMessage_TypeNegativeNoDecimalSignNoWallet(t *testing.T) {
	// Given
	msg := "$-123 Negative No Decimal Sign No Wallet No Category"
	ctrl := NewParserController(nil, nil, nil, nil)

	// When
	msgType := ctrl.GetTypeFromMessage(&tg.Message{Payload: msg})

	// Then
	require.EqualValues(t, messageTypeTransaction, msgType)
}
func TestParserController_GetTypeFromMessage_Unknown(t *testing.T) {
	// Given
	msg := "abcdefg"
	ctrl := NewParserController(nil, nil, nil, nil)

	// When
	msgType := ctrl.GetTypeFromMessage(&tg.Message{Payload: msg})

	// Then
	require.EqualValues(t, messageTypeUnknown, msgType)
}
