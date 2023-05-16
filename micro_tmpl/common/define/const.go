package define

const (
	GET  = 1
	POST = 2
)

const (
	Success                 = 200
	ParamError              = 401
	TokenError              = 402
	PathError               = 403
	NotFoundLeaderboard     = 404
	CookieError             = 405
	ServerError             = 500
	Timeout                 = 503
	ServerPause             = 504
	UserNameOrPasswordError = 505
)

const (
	ParamErrorMsg              = "ParamError"
	TokenErrorMsg              = "TokenInvalid"
	PathErrorMsg               = "RequestPathError"
	NotFoundLeaderboardMsg     = "NotFoundLeaderboard"
	CookieErrorMsg             = "CookieError"
	TimeoutMsg                 = "Timeout"
	ServerErrorMsg             = "ServerError"
	ServerPauseMsg             = "ServerPause"
	UserNameOrPasswordErrorMsg = "UserNameOrPasswordError"
)
