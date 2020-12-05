package creditVoucher

import (
	"encoding/json"
	"errors"
	"r1wallet/models"
	"r1wallet/repositories"
	"strconv"
)

type CreditVoucher struct {
	repository     *repositories.Repository
	publishChannel string
}

var VoucherSoldOut = errors.New("voucher sold out")

func NewCreditVoucher(repository *repositories.Repository, publishChannel string) *CreditVoucher {
	return &CreditVoucher{
		repository:     repository,
		publishChannel: publishChannel,
	}
}

func (c *CreditVoucher) Redeem(userID int, code string) error {
	voucher, err := c.repository.Voucher.FindVoucherByCode(code)
	if err != nil {
		return err
	}

	v, err := c.getUsedCount(code)
	if err != nil {
		return err
	}

	if voucher.Usable <= v {
		return VoucherSoldOut
	}

	return c.repository.Voucher.RedeemVoucher(userID, voucher, c.sendIncreaseRequestToWallet)

}

func (c *CreditVoucher) getUsedCount(code string) (int, error) {
	var v int
	val, err := c.repository.Redis.GetValue(code)
	if err != nil {
		v, err = c.repository.Voucher.GetRedeemedCount(code)
		if err != nil {
			return 0, err
		}

		_ = c.repository.Redis.SetValue(code, v)
	} else {
		v, err = strconv.Atoi(val)
		if err != nil {
			return 0, err
		}
	}
	return v, nil
}

func (c *CreditVoucher) sendIncreaseRequestToWallet(userID int, voucher models.VoucherModel) error {
	d, err := json.Marshal(models.IncreaseRequestModel{
		UserID: userID,
		Amount: voucher.Amount,
	})
	if err != nil {
		return err
	}

	err = c.repository.Redis.Increase(voucher.Code)
	if err != nil {
		return err
	}

	err = c.repository.Redis.PublishInChannel(d, c.publishChannel)
	if err != nil {
		return err
	}

	return nil
}
