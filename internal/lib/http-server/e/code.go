package e

const (
	// Err0System System error.
	Err0System = iota
)

const (
	// 401 Unauthorized errors.
	_ = 40100000 + iota
	Err401AuthEmptyTokenError
	Err401TokenError
	Err401UserNotFoundError
	Err401UserNotActiveError
	Err401SystemEmptyTokenError
	Err401SystemTokenError
)

const (
	// 404 Not Found errors.
	_ = 40400000 + iota
	Err404URLNotFound
	Err404URLExpired
)

const (
	// 422 Unprocessable Entity errors.
	_ = 42200000 + iota
	Err422URLValidateCreateError
	Err422URLSaveError
	Err422URLFindError
	Err422LoginUserNotFoundError
	Err422LoginAuthTokensError
	Err422LoginComparePasswordError = 42200024
	Err422LoginInvalidPasswordError = 42200025
)
