package user

const (
	CreatePending = "CREATE_PENDING"
	Created = "CREATED"
	Rejected = "REJECTED"
	Remove = "REMOVED"
)

const (
	DuplicateErrorCode = 1062
	DataTooLongErrorCode = 1406
)

const (
	KeyUserId = "uix_auths_user_id"
	KeyMsgId = "uix_processed_messages_msg_id"
)

const (
	ColumnUserId = "user_id"
	ColumnMsgId = "msg_id"
)