package response

const (
	SUCCESS   = 200
	FAIL      = 400
	AuthError = 401
	ERROR     = 500

	UserExist       = 1001
	UserNotExist    = 1002
	PasswordError   = 1003
	ParamError      = 1004
	PermissionError = 1005
)
