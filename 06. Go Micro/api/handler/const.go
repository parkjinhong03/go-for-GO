package handler

import "time"

const (
	DefaultDialTimeout = time.Second * 2
	DefaultRequestTimeout = time.Second * 3
)

const (
	emailDuplicateIndex = 0
	userIdDuplicateIndex = 0
	userCreateIndex = 1
)