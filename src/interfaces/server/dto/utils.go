// Package dto
package dto

import (
	"github.com/labstack/echo/v4"
)

var (
	ErrLackParam             = NewApiStatus("PARAM_MISS", "缺少参数", HttpCodeBadRequest)
	ErrInvalidParam          = NewApiStatus("PARAM_INVALID", "非法参数", HttpCodeBadRequest)
	ErrErrorParam            = NewApiStatus("PARAM_ERROR", "参数错误", HttpCodeBadRequest)
	ErrServerError           = NewApiStatus("SERVER_ERROR", "服务器错误", HttpCodeInternalError)
	ErrRateLimitExceeded     = NewApiStatus("RATE_LIMIT_EXCEEDED", "请求频率过高", HttpCodeTooManyRequests)
	ErrNoMatchRoute          = NewApiStatus("NO_MATCH_ROUTE", "未匹配到路由", HttpCodeNotFound)
	SuccessHandleRequest     = NewApiStatus("SUCCESS", "成功", HttpCodeOk)
	ErrMissingOrMalformedJwt = NewApiStatus("MISSING_OR_MALFORMED_JWT", "缺少JWT令牌或者令牌格式错误", HttpCodeBadRequest)
	ErrInvalidOrExpiredJwt   = NewApiStatus("INVALID_OR_EXPIRED_JWT", "无效或过期的JWT令牌", HttpCodeUnauthorized)
	ErrInvalidJwtType        = NewApiStatus("INVALID_JWT_TYPE", "非法的JWT令牌类型", HttpCodeUnauthorized)
	ErrUnknownJwtError       = NewApiStatus("UNKNOWN_JWT_ERROR", "未知的JWT解析错误", HttpCodeInternalError)
)

func ErrorResponse(ctx echo.Context, codeStatus *ApiStatus) error {
	return NewApiResponse[any](codeStatus, nil).Response(ctx)
}

func TextResponse(ctx echo.Context, httpCode int, content string) error {
	return ctx.String(httpCode, content)
}
