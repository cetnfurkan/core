package http

import "github.com/labstack/echo/v4"

type response struct {
	*Pagination
	Code string `json:"code"`
}

func Response(ctx echo.Context, code string, data any, pagination *Pagination) error {
	if pagination != nil {
		pagination.Data = data
	} else {
		pagination = &Pagination{
			Data: data,
		}
	}

	resp := response{
		Code:       code,
		Pagination: pagination,
	}

	return ctx.JSON(200, resp)
}
