package models

import "github.com/pdragnev/notification-system/common"

type RetryError struct {
	Msg            string
	UpdatedMessage common.NotificationMessage
}

func (e *RetryError) Error() string {
	return e.Msg
}

func NewRetryError(msg string, updatedMsg common.NotificationMessage) error {
	return &RetryError{
		Msg:            msg,
		UpdatedMessage: updatedMsg,
	}
}

type DeserializingMsgError struct {
	Msg string
}

func (e *DeserializingMsgError) Error() string {
	return e.Msg
}

func NewDeserializingMsgError(msg string) error {
	return &DeserializingMsgError{
		Msg: msg,
	}
}

type ProcessingTypeError struct {
	Msg string
}

func (e *ProcessingTypeError) Error() string {
	return e.Msg
}

func NewProcessingTypeError(msg string) error {
	return &ProcessingTypeError{
		Msg: msg,
	}
}
