package common

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type CustomContext struct {
	echo.Context
}

func (c *CustomContext) BindAndValidate(i interface{}) (err error) {
	if err = c.Bind(i); err != nil {
		return
	}
	err = c.Validate(i)
	return
}

func (c *CustomContext) OkPagination(i interface{}, params PaginateParams) (err error) {
	return c.JSON(http.StatusOK, NewSlicePaginator(i, params.Page, params.Limit).Json())
}
