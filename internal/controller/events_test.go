package controller

import (
	tg "gopkg.in/telebot.v3"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParserController_GetTypeFromMessage_TypeNoDecimalNoSignNoWallet(t *testing.T) {
	// Given
	msg := "123 No Decimal No Sign No Category No Wallet"
	ctrl := NewEventsController(nil, nil, nil, nil, nil)

	// When
	msgType := ctrl.GetTypeFromMessage(&tg.Message{
		Chat: &tg.Chat{Type: tg.ChatPrivate},
		Text: msg,
	})

	// Then
	require.EqualValues(t, messageTypeTransaction, msgType)
}
func TestParserController_GetTypeFromMessage_TypeDotDecimalNoSignNoWallet(t *testing.T) {
	// Given
	msg := "123.45 Dot Decimal No Sign No Category No Wallet"
	ctrl := NewEventsController(nil, nil, nil, nil, nil)

	// When
	msgType := ctrl.GetTypeFromMessage(&tg.Message{
		Chat: &tg.Chat{Type: tg.ChatPrivate},
		Text: msg,
	})

	// Then
	require.EqualValues(t, messageTypeTransaction, msgType)
}
func TestParserController_GetTypeFromMessage_TypeCommaDecimalNoSignNoWallet(t *testing.T) {
	// Given
	msg := "123,45 Comma Decimal No Sign No Category No Wallet"
	ctrl := NewEventsController(nil, nil, nil, nil, nil)

	// When
	msgType := ctrl.GetTypeFromMessage(&tg.Message{
		Chat: &tg.Chat{Type: tg.ChatPrivate},
		Text: msg,
	})

	// Then
	require.EqualValues(t, messageTypeTransaction, msgType)
}
func TestParserController_GetTypeFromMessage_TypeNoDecimalSignNoWallet(t *testing.T) {
	// Given
	msg := "$123 No Decimal No Sign No Category No Wallet"
	ctrl := NewEventsController(nil, nil, nil, nil, nil)

	// When
	msgType := ctrl.GetTypeFromMessage(&tg.Message{
		Chat: &tg.Chat{Type: tg.ChatPrivate},
		Text: msg,
	})

	// Then
	require.EqualValues(t, messageTypeTransaction, msgType)
}
func TestParserController_GetTypeFromMessage_TypeDotDecimalSignNoWallet(t *testing.T) {
	// Given
	msg := "$123.45 Dot Decimal No Sign No Category No Wallet"
	ctrl := NewEventsController(nil, nil, nil, nil, nil)

	// When
	msgType := ctrl.GetTypeFromMessage(&tg.Message{
		Chat: &tg.Chat{Type: tg.ChatPrivate},
		Text: msg,
	})

	// Then
	require.EqualValues(t, messageTypeTransaction, msgType)
}
func TestParserController_GetTypeFromMessage_TypeCommaDecimalSignNoWallet(t *testing.T) {
	// Given
	msg := "$123,45 Comma Decimal No Sign No Category No Wallet"
	ctrl := NewEventsController(nil, nil, nil, nil, nil)

	// When
	msgType := ctrl.GetTypeFromMessage(&tg.Message{
		Chat: &tg.Chat{Type: tg.ChatPrivate},
		Text: msg,
	})

	// Then
	require.EqualValues(t, messageTypeTransaction, msgType)
}
func TestParserController_GetTypeFromMessage_TypeNegativeNoDecimalNoSignNoWallet(t *testing.T) {
	// Given
	msg := "-123 Negative No Decimal No Sign No Wallet No Category"
	ctrl := NewEventsController(nil, nil, nil, nil, nil)

	// When
	msgType := ctrl.GetTypeFromMessage(&tg.Message{
		Chat: &tg.Chat{Type: tg.ChatPrivate},
		Text: msg,
	})

	// Then
	require.EqualValues(t, messageTypeTransaction, msgType)
}
func TestParserController_GetTypeFromMessage_TypeNegativeNoDecimalSignNoWallet(t *testing.T) {
	// Given
	msg := "$-123 Negative No Decimal Sign No Wallet No Category"
	ctrl := NewEventsController(nil, nil, nil, nil, nil)

	// When
	msgType := ctrl.GetTypeFromMessage(&tg.Message{
		Chat: &tg.Chat{Type: tg.ChatPrivate},
		Text: msg,
	})

	// Then
	require.EqualValues(t, messageTypeTransaction, msgType)
}
func TestParserController_GetTypeFromMessage_Unknown(t *testing.T) {
	// Given
	msg := "abcdefg"
	ctrl := NewEventsController(nil, nil, nil, nil, nil)

	// When
	msgType := ctrl.GetTypeFromMessage(&tg.Message{
		Chat: &tg.Chat{Type: tg.ChatPrivate},
		Text: msg,
	})

	// Then
	require.EqualValues(t, messageTypeUnknown, msgType)
}
