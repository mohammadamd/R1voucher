package handler

import (
	"github.com/labstack/echo/v4"
	"r1wallet/services"
)

type BaseHandler struct {
	Voucher Voucher
}

type Voucher interface {
	RedeemVoucher() func(c echo.Context) error
}

func NewBaseHandler(services *services.Services) *BaseHandler {
	return &BaseHandler{
		Voucher: newVoucherHandler(services),
	}
}
