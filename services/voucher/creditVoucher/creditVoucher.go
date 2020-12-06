package creditVoucher

import (
	"encoding/json"
	"errors"
	"fmt"
	"r1wallet/models"
	"r1wallet/repositories"
	"strconv"
)

type CreditVoucher struct {
	repository             *repositories.Repository
	communicationQueueName string
}

var VoucherSoldOut = errors.New("voucher sold out")

func NewCreditVoucher(repository *repositories.Repository, comQueue string) *CreditVoucher {
	return &CreditVoucher{
		repository:             repository,
		communicationQueueName: comQueue,
	}
}

func (c *CreditVoucher) Redeem(userID int, code string) error {
	voucher, err := c.repository.Voucher.FindVoucherByCode(code)
	if err != nil {
		return err
	}

	v, err := c.getUsedCount(voucher.ID)
	if err != nil {
		return err
	}

	if voucher.Usable <= v {
		return VoucherSoldOut
	}

	return c.repository.Voucher.RedeemVoucher(userID, voucher, c.getStep, c.sendIncreaseRequestToWallet)

}

func (c *CreditVoucher) getUsedCount(voucherID int) (int, error) {
	var v int
	val, err := c.repository.Redis.GetValue(getRedisCacheKeyForVoucher(voucherID))
	if err != nil {
		v, err = c.repository.Voucher.GetRedeemedCount(voucherID)
		if err != nil {
			return 0, err
		}

		_ = c.repository.Redis.SetValue(getRedisCacheKeyForVoucher(voucherID), v)
	} else {
		v, err = strconv.Atoi(val)
		if err != nil {
			return 0, err
		}
	}
	return v, nil
}

func (c *CreditVoucher) getStep(voucher models.VoucherModel) (int, error) {
	cv, err := c.repository.Redis.Increase(getRedisCacheKeyForVoucher(voucher.ID))
	if err != nil {
		return 0, err
	}

	if cv > voucher.Usable {
		return 0, VoucherSoldOut
	}

	return cv, nil
}

func (c *CreditVoucher) sendIncreaseRequestToWallet(userID int, voucher models.VoucherModel) error {
	d, err := json.Marshal(models.IncreaseRequestModel{
		UserID: userID,
		Amount: voucher.Amount,
	})
	if err != nil {
		return err
	}

	err = c.repository.Redis.Enqueue(d, c.communicationQueueName)
	if err != nil {
		return err
	}

	return nil
}

func getRedisCacheKeyForVoucher(voucherID int) string {
	return fmt.Sprintf("voucher:%d", voucherID)
}
