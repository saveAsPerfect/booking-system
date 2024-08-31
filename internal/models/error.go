package models

import "errors"

var (
	ErrorRoomAlreadyReserved = errors.New("the room is reserved for this time")
	ErrorInvalidData = errors.New("invalid data")
)
