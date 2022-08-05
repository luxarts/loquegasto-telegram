package defines

import "github.com/luxarts/jsend-go"

var (
	ErrInvalidID          = jsend.NewFail("invalid id")
	ErrInvalidMsgID       = jsend.NewFail("invalid msgID")
	ErrInvalidBody        = jsend.NewFail("invalid body")
	ErrNameAlreadyExists  = jsend.NewFail(map[string]string{"name": "already used"})
	ErrEmojiAlreadyExists = jsend.NewFail(map[string]string{"emoji": "already used"})
	ErrUserAlreadyExists  = jsend.NewFail(map[string]string{"user": "already used"})
)
