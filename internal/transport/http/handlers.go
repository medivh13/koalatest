package http

import (
	"net/http"
	"os"

	"github.com/apex/log"
	"github.com/dchest/uniuri"
	"github.com/medivh13/koalatest/internal/services"
	koalaConst "github.com/medivh13/koalatest/pkg/common/const"
	"github.com/medivh13/koalatest/pkg/common/crypto"
	"github.com/medivh13/koalatest/pkg/dto"
	btbErrors "github.com/medivh13/koalatest/pkg/errors"

	"github.com/labstack/echo"
)

type HttpHandler struct {
	service services.Services
}

func NewHttpHandler(e *echo.Echo, srv services.Services) {
	handler := &HttpHandler{
		srv,
	}

	e.GET("api/koalatest/ping", handler.Ping)
	e.GET("api/koalatest/register", handler.PostRegister)

}

func (h *HttpHandler) Ping(c echo.Context) error {

	version := os.Getenv("VERSION")
	if version == "" {
		version = "pong"
	}

	data := version

	return c.JSON(http.StatusOK, data)

}

func (h *HttpHandler) PostRegister(c echo.Context) error {

	postDTO := dto.CustomersReqDTO{}

	if err := c.Bind(&postDTO); err != nil {
		log.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}
	postDTO.CustomerId = uniuri.NewLen(15)
	postDTO.Salt = uniuri.NewLen(10)
	postDTO.Password = crypto.EncodeSHA256HMAC(postDTO.Salt, postDTO.Password)

	err := h.service.Register(&postDTO)

	if err != nil {
		log.Error(err.Error())
		return c.JSON(getStatusCode(err), map[string]string{
			"error": err.Error(),
		})
	}

	var resp = dto.ResponseDTO{
		Success: true,
		Message: koalaConst.SaveSuccess,
		Data:    nil,
	}

	return c.JSON(http.StatusOK, resp)

}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case btbErrors.ErrInternalServerError:
		return http.StatusInternalServerError
	case btbErrors.ErrNotFound:
		return http.StatusNotFound
	case btbErrors.ErrConflict:
		return http.StatusConflict
	case btbErrors.ErrInvalidRequest:
		return http.StatusBadRequest
	case btbErrors.ErrFailAuth:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
