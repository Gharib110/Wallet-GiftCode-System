package repo

import (
	"context"
	"github.com/Gharib110/ArvanWallet/models"
	"github.com/Gharib110/ArvanWallet/queries"
	"time"
)

type Transactions interface {
	BeginTransaction(context.Context, *models.TransactionDTO) (error, *models.TransactionDTO)
}

func (c *DBConfig) BeginTransaction(ctx context.Context,
	dto *models.TransactionDTO) (error, *models.TransactionDTO) {
	ctxx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	tx, err := c.dB.BeginTx(ctxx, nil)

	if err != nil {
		c.errorLogger.Println(err)
		tx.Rollback()

		dto.Type = "FAILED"
		return err, dto
	}

	t := time.Now()
	_, err = tx.Exec(queries.InserIntoTransactionsTable, &dto.UserID, &dto.GiftCodeID,
		&dto.Amount, "CREATED-DEPOSIT", t)
	if err != nil {
		c.errorLogger.Println(err)
		tx.Rollback()

		dto.Timestamp = time.Now()
		dto.Type = "FAILED"
		return err, dto
	}

	userDao := &models.UserDAO{
		ID:     dto.UserID,
		Charge: int64(dto.UserCharge + dto.Amount),
	}

	err = c.UpdateUserCharge(ctxx, userDao)

	if err != nil {
		t = time.Now()
		_, err = tx.Exec(queries.InserIntoTransactionsTable, &dto.UserID,
			&dto.GiftCodeID, &dto.Amount, "FAILED-DEPOSIT", t)
		c.errorLogger.Println(err)
		tx.Rollback()

		dto.Timestamp = t
		dto.Type = "FAILED"
		return err, dto
	}

	t = time.Now()
	_, err = tx.Exec(queries.InserIntoTransactionsTable, &dto.UserID,
		&dto.GiftCodeID, &dto.Amount, "DONE-DEPOSIT", t)
	if err != nil {
		c.errorLogger.Println(err)
		tx.Rollback()

		dto.Timestamp = t
		dto.Type = "FAILED"
		return err, dto
	}

	err = tx.Commit()
	if err != nil {
		c.errorLogger.Println(err)
		return err, nil
	}

	t = time.Now()
	dto.Timestamp = t
	dto.Type = "DONE"
	return nil, dto
}
