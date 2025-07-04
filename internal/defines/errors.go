package defines

import (
	"errors"
	"github.com/luxarts/jsend-go"
)

var (
	// API errors
	ErrInvalidID          = jsend.NewFail("invalid id")
	ErrInvalidMsgID       = jsend.NewFail("invalid msgID")
	ErrInvalidBody        = jsend.NewFail("invalid body")
	ErrNameAlreadyExists  = jsend.NewFail(map[string]string{"name": "already used"})
	ErrEmojiAlreadyExists = jsend.NewFail(map[string]string{"emoji": "already used"})
	ErrUserAlreadyExists  = jsend.NewFail(map[string]string{"user": "already used"})

	// Internal errors
	ErrInvalidSyntax = errors.New("invalid syntax")

	ErrSessionNotFound = errors.New("session not found")
)
