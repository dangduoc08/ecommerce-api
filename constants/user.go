package constants

const (
	USER_STATUS_FIELD_NAME = "user_status"
)

const (
	USER_ACTIVE    = "active"
	USER_INACTIVE  = "inactive"
	USER_SUSPENDED = "suspended"
)

var USER_STATUSES = []string{
	USER_ACTIVE,
	USER_INACTIVE,
	USER_SUSPENDED,
}
