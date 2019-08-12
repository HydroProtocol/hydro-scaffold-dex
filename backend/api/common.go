package api

import (
	"github.com/labstack/echo"
	v "gopkg.in/go-playground/validator.v9"
	"reflect"
)

var validate = v.New()

func bindAndValidParams(c echo.Context, params Param) (err error) {
	cc := c.(*HydroApiContext)
	params.SetAddress(cc.Address)

	if err := c.Bind(params); err != nil {
		panic(BindError())
	}

	bindUrlParam(c, params)
	return validate.Struct(params)
}

func bindUrlParam(c echo.Context, ptr interface{}) {
	typ := reflect.TypeOf(ptr).Elem()
	val := reflect.ValueOf(ptr).Elem()

	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := val.Field(i)
		if !structField.CanSet() {
			continue
		}
		inputFieldName := typeField.Tag.Get("param")
		if inputFieldName == "" {
			continue
		}

		structField.SetString(c.Param(inputFieldName))
		continue
	}
}

func commonHandler(params Param, fn func(Param) (interface{}, error)) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		var req Param
		if params != nil {
			// Params is shared among all request
			// We create a new one for each request
			req = reflect.New(reflect.TypeOf(params).Elem()).Interface().(Param)
			err = bindAndValidParams(c, req)
		}

		if err != nil {
			return
		}

		resp, err := fn(req)

		if err != nil {
			return
		}

		return commonResponse(c, resp)
	}
}
