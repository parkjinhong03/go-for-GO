package handler

const (
	StatusEmailDuplicated = 471
)

const (
	MessageEmailNotDuplicated = "you can use this email"
	MessageEmailDuplicated = "this email is already in use"
)

const (
	MessageBadRequest = "the request is invalid. please check the document"
)

const (
	MessageUnableGetMetadata = "unable to get metadata from context"
	MessageThereIsNoXReqId = "X-Request-Id is not included in message header"
	MessageNoSpanContext = "there ins't span context in metadata"
)