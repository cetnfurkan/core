package http

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

type Pagination struct {
	Data         any     `json:"data,omitempty"`
	Search       *string `json:"search,omitempty"`
	TotalElement *int    `json:"total_element,omitempty"`
	Page         *int    `json:"page,omitempty"`
	Size         *int    `json:"size,omitempty"`
}

func NewPagination(search string, page, size int) *Pagination {
	return &Pagination{
		Search: &search,
		Page:   &page,
		Size:   &size,
	}
}

func Paginate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			search := c.QueryParam("search")
			p, _ := strconv.Atoi(c.QueryParam("page"))
			s, _ := strconv.Atoi(c.QueryParam("size"))

			pg := NewPagination(search, p, s)
			c.Set("page", pg)

			return next(c)
		}
	}
}

func GetPagination(c echo.Context) *Pagination {
	p := c.Get("page")

	page, ok := p.(*Pagination)
	if !ok {
		return &Pagination{}
	}

	return page
}

func (p *Pagination) Offset() int {
	if p.Page == nil || p.Size == nil {
		return 0
	}

	return (*(p.Page) - 1) * *(p.Size)
}
