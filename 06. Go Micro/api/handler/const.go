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

const (
	apiGateway = "examples.blog.api.gateWay"
	authClient = "examples.blog.client.auth"
	userClient = "examples.blog.client.user"
)