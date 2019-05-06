package api

import (
	"github.com/labstack/echo"
	"strings"
)

type HydroApiContext struct {
	echo.Context
	// If address is not empty means this user is authenticated.
	Address string
}

func initHydroApiContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := &HydroApiContext{c, ""}
		return next(cc)
	}
}

func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(*HydroApiContext)
		cc.Response().Header().Set(echo.HeaderServer, "Echo/3.0")

		hydroAuthToken := cc.Request().Header.Get("Hydro-Authentication")
		hydroAuthTokens := strings.Split(hydroAuthToken, "#")

		if len(hydroAuthTokens) != 3 {
			return &ApiError{Code: -11, Desc: "Hydro-Authentication should be like {address}#HYDRO-AUTHENTICATION@{time}#{signature}"}
		}

		valid, err := hydro.IsValidSignature(hydroAuthTokens[0], hydroAuthTokens[1], hydroAuthTokens[2])
		if !valid || err != nil {
			return &ApiError{Code: -11, Desc: "Hydro-Authentication valid failed, please check your authentication"}
		}
		cc.Address = strings.ToLower(hydroAuthTokens[0])
		return next(cc)
	}
}
