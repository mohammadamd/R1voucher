package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"r1wallet/models"
	"r1wallet/repositories"
	"r1wallet/services"
	"r1wallet/services/voucher/creditVoucher"
)

type VoucherHandler struct {
	service *services.Services
}

func newVoucherHandler(service *services.Services) *VoucherHandler {
	return &VoucherHandler{service: service}
}

func (vh *VoucherHandler) RedeemVoucher() func(c echo.Context) error {
	return func(c echo.Context) error {
		var rq models.RedeemVoucherRequest
		err := c.Bind(&rq)
		if err != nil {
			fmt.Println("could not bind redeem request")
			return err
		}

		if rq.Code == "" || rq.UserID == 0 {
			return echo.ErrBadRequest
		}

		err = vh.service.Voucher.Redeem(rq.UserID, rq.Code)
		switch err {
		case creditVoucher.VoucherSoldOut:
			return c.JSON(http.StatusNotAcceptable, map[string]interface{}{"message": "Sorry voucher code sold out! :("})
		case repositories.InvalidVoucherCode:
			return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"message": "Entered voucher code is invalid"})
		case repositories.VoucherAlreadyUsed:
			return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"message": "You have already used this code"})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"message": "Congratulation your credit will be added to your wallet soon"})
	}
}
