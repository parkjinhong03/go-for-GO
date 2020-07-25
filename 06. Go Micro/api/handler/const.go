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
	authClient = "examples.blog.client.auth"
	userClient = "examples.blog.client.user"
)

const (
	userIdDuplicate = "examples.blog.service.auth.UserIdDuplicate"
	beforeCreateAuth = "examples.blog.service.auth.BeforeCreateAuth"
	emailDuplicate = "examples.blog.service.user.EmailDuplicate"
)
