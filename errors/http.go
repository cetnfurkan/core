package errors

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type (
	HttpErrorCode string

	httpErrorInfo struct {
		Code         HttpErrorCode
		Message      string
		ResponseCode int
		ServiceName  string
	}

	httpError struct {
		ctx echo.Context
		httpErrorInfo
	}

	httpErrorOption func(*httpError)
)

func HttpErrorInfo(code HttpErrorCode, message string, responseCode int) *httpErrorInfo {
	return &httpErrorInfo{
		Code:         code,
		Message:      message,
		ResponseCode: responseCode,
	}
}

func HttpError(ctx echo.Context, errorInfo *httpErrorInfo, opts ...httpErrorOption) error {
	err := &httpError{
		ctx: ctx,
	}

	if errorInfo != nil {
		err.httpErrorInfo = *errorInfo
	}

	for _, opt := range opts {
		opt(err)
	}

	return ctx.JSON(err.ResponseCode, map[string]any{
		"code":    err.Code,
		"message": err.Message,
		"service": err.ServiceName,
	})
}

func WithErrorCode(code HttpErrorCode) httpErrorOption {
	return func(err *httpError) {
		err.Code = code
	}
}

func WithErrorMessage(message string) httpErrorOption {
	return func(err *httpError) {
		err.Message = message
	}
}

func WithResponseCode(code int) httpErrorOption {
	return func(err *httpError) {
		err.ResponseCode = code
	}
}

func WithServiceName(name string) httpErrorOption {
	return func(err *httpError) {
		err.ServiceName = name
	}
}

func (err *httpError) Error() string {
	return fmt.Sprintf("[%s]-[%s]: %s", err.ServiceName, err.Code, err.Message)
}
