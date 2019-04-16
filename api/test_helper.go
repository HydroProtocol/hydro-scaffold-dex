package api

import (
	"encoding/json"
	"github.com/labstack/echo"
	"io"
	"net/http/httptest"
	"strings"
)

func request(url, method, auth string, body interface{}) *Response {
	e := getEchoServer()
	var reader io.Reader
	if body == nil {
		reader = nil
	} else {
		bts, _ := json.Marshal(body)
		strings.NewReader(string(bts))
	}

	req := httptest.NewRequest(method, url, reader)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	if auth == "" {
		address := "0x5409ed021d9299bf6814279a6a1411a7e866a631"
		signature := "0xdcd19ecc53c51bc1c8c67183d9ed8a2c68bb3717b7bbbd39da969960feeb95d45f79ead1d476c5cb1f2ebf77b76a87abee2bf5643a235125a85428d3ef4926b700"
		message := "HYDRO-AUTHENTICATION"
		auth = address + "#" + message + "#" + signature
	}

	req.Header.Set("Hydro-Authentication", auth)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	var res Response
	json.Unmarshal(rec.Body.Bytes(), &res)
	return &res
}
