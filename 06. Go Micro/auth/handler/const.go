package handler

const (
	StatusUserIdDuplicate = 470
	StatusEmailDuplicate = 471
)

// BeforeCreateAuth 메서드에서만 사용하는 Message
const (
	MessageAuthCreated = "user created reservation has been successfully processed"
	MessageBadRequest = "the request is not valid. please check the document"
)

// UserIdDuplicated 메서드에서만 사용하는 Message
const (
	MessageUserIdNotDuplicated = "you can use that user id"
)

const (
	MessageUserIdDuplicate = "this user id is already in use"
	MessageEmailDuplicate = "this email is already in use"
)