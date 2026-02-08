package service

import (
	"errors"
	"server-a/server/constant/message"
)

var (
	errAppleSignInFailed = errors.New(message.AppleSignInFailed)
)
