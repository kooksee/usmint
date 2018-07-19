package web

import (
	"github.com/labstack/echo"
	"net/http"
	"github.com/kooksee/usmint/types"
)

func txPost(c echo.Context) error {
	tx := types.NewTransaction()

	if err := c.Bind(tx); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(tx); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, m{
		"ok": "ok",
		"oo": m{
			"rid": c.Response().Header().Get(echo.HeaderXRequestID),
		},
		"d": tx,
	})
}

func txGet(c echo.Context) error {
	txId := c.Param("id")
	return c.JSON(http.StatusOK, txId)
}
