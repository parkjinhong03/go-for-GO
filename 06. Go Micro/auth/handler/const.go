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
	MessageUnableGetMetadata = "unable to get metadata from context"
	MessageThereIsNoXReqId = "X-Request-Id is not included in message header"
	MessageThereIsNoCert = "Unique-Authorization is not included in message header"
	MessageUserIdDuplicate = "this user id is already in use"
	MessageEmailDuplicate = "this email is already in use"
	MessageNoSpanContext = "there ins't span context in metadata"
)