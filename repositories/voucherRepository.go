package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"r1wallet/models"
)

type voucherRepository struct {
	db  *sql.DB
	dbq dbQE
	tx  *sql.Tx
}

type dbQE interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

const (
	validateVoucherQuery         = "SELECT (id, code, usable, amount) FROM vouchers WHERE code = ? limit 1"
	redeemVoucherQuery           = "INSERT INTO redeemed_voucher (user_id, voucher_id) VALUES (?, ?)"
	redeemVoucherCountQuery      = "SELECT COUNT(1) from redeemed_voucher WHERE code = ?"
	validateFirstTimeRedeemQuery = "SELECT 1 from redeemed_voucher WHERE voucher_id = ? AND user_id = ?"
)

var InvalidVoucherCode = errors.New("invalid voucher code")
var VoucherAlreadyUsed = errors.New("voucher already redeemed by user")

func NewVoucherRepository(db *sql.DB) *voucherRepository {
	return &voucherRepository{
		db:  db,
		dbq: db,
		tx:  nil,
	}
}

func (v *voucherRepository) FindVoucherByCode(code string) (models.VoucherModel, error) {
	var vm models.VoucherModel
	rows, err := v.dbq.Query(validateVoucherQuery, code)
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println("failed to close rows: ", err)
		}
	}(rows)

	if err != nil {
		return vm, err
	}

	if rows.Next() == false {
		return vm, InvalidVoucherCode
	}

	err = rows.Scan(&vm.ID, &vm.Code, &vm.Usable, &vm.Amount)
	if err != nil {
		return vm, err
	}

	return vm, nil
}

func (v *voucherRepository) InsertIntoRedeemedVoucher(userID, voucherID int) error {
	rows, err := v.dbq.Exec(redeemVoucherQuery, userID, voucherID)
	if err != nil {
		return err
	}

	ra, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if ra == 0 {
		return fmt.Errorf("failed to insert into redeemed voucher, values : %d, %d", userID, voucherID)
	}

	return nil
}

func (v *voucherRepository) IsUserRedeemedVoucherBefore(userID, voucherID int) (bool, error) {
	rows, err := v.dbq.Query(validateFirstTimeRedeemQuery, voucherID, userID)
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println("failed to close rows: ", err)
		}
	}(rows)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return true, err
	}

	return true, nil
}

func (v *voucherRepository) GetRedeemedCount(code string) (int, error) {
	rows, err := v.dbq.Query(redeemVoucherCountQuery, code)
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println("failed to close rows: ", err)
		}
	}(rows)

	if err != nil {
		return 0, err
	}

	var i int
	err = rows.Scan(&i)
	if err != nil {
		return 0, err
	}

	return i, nil
}

func (v *voucherRepository) RedeemVoucher(userID int, voucher models.VoucherModel, closure func(userID int, voucher models.VoucherModel) error) error {
	trx, err := v.beginTransaction()
	if err != nil {
		return err
	}

	ub, err := trx.IsUserRedeemedVoucherBefore(userID, voucher.ID)
	if ub {
		er := trx.rollbackTransaction()
		if er != nil {
			return er
		}

		return VoucherAlreadyUsed
	}

	err = trx.InsertIntoRedeemedVoucher(userID, voucher.ID)
	if err != nil {
		er := trx.rollbackTransaction()
		if er != nil {
			return er
		}

		return err
	}

	err = closure(userID, voucher)
	if err != nil {
		er := trx.rollbackTransaction()
		if er != nil {
			return er
		}

		return err
	}

	return trx.commitTransaction()
}


func (v *voucherRepository) beginTransaction() (*voucherRepository, error) {
	tx, err := v.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return &voucherRepository{}, err
	}

	return &voucherRepository{tx: tx, dbq: tx}, nil
}

func (v *voucherRepository) commitTransaction() error {
	if v.tx == nil {
		return fmt.Errorf("you cant commit tansaction befor start it")
	}

	return v.tx.Commit()
}

func (v *voucherRepository) rollbackTransaction() error {
	if v.tx == nil {
		return fmt.Errorf("you cant rollback tansaction befor start it")
	}

	return v.tx.Rollback()
}
