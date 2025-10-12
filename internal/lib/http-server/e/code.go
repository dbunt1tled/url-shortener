package e

const (
	Err422URLValidateCreateError    = 42200001
	Err422URLSaveError              = 42200002
	Err422URLFindError              = 42200003
	Err422LoginUserNotFoundError    = 42200004
	Err422LoginAuthTokensError      = 42200005
	Err422LoginComparePasswordError = 42200024
	Err422LoginInvalidPasswordError = 42200025

	Err401AuthEmptyTokenError   = 40100001
	Err401TokenError            = 40100002
	Err401UserNotFoundError     = 40100003
	Err401UserNotActiveError    = 40100004
	Err401SystemEmptyTokenError = 40100005
	Err401SystemTokenError      = 40100005

	Err404URLNotFound = 40100001
	Err404URLExpired  = 40100002
)
