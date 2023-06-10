package repo

import (
	"context"
	"github.com/Gharib110/ArvanGift/models"
	"github.com/Gharib110/ArvanGift/queries"
	"time"
)

type GiftCode interface {
	RegisterGiftCode(context.Context, *models2.GiftCodeDAO) error
	GetGiftCodeByID(context.Context, *models2.GiftCodeDAO) *models2.GiftCode
	UpdateGiftCodeCount(context.Context, *models2.GiftCodeDAO) error
}

func (c *DBConfig) RegisterGiftCode(ctx context.Context, gift *models2.GiftCodeDAO) error {
	PingPQDB()
	ctxx, cancel := context.WithTimeout(ctx, time.Minute*1) // Time should not be hardCoded
	defer cancel()

	result, err := c.dB.ExecContext(ctxx, queries2.InsertIntoGiftCodeTable,
		&gift.Code, &gift.Amount, &gift.IsActive, &gift.RedemptionLimit,
		&gift.RedemptionCount, &gift.StartTime, &gift.ExpirationTime)

	if err != nil {
		c.errorLogger.Println(err)
		return err
	}

	gift.ID, _ = result.LastInsertId()
	c.infoLogger.Println("GiftCode is Created !")
	return nil
}

func (c *DBConfig) GetGiftCodeByID(ctx context.Context,
	gift *models2.GiftCodeDAO) *models2.GiftCode {
	PingPQDB()
	ctxx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	row := c.dB.QueryRowContext(ctxx, queries2.GetGiftCodeByCodeID, &gift.Code)
	if row.Err() != nil {
		c.errorLogger.Println(row.Err())
		return nil
	}

	g := &models2.GiftCode{
		ID:              0,
		Code:            "",
		Amount:          0,
		IsActive:        false,
		RedemptionLimit: 0,
		RedemptionCount: 0,
		StartTime:       "",
		ExpirationTime:  "",
	}
	st := time.Time{}
	exp := time.Time{}

	err := row.Scan(&g.ID, &g.Code, &g.Amount, &g.IsActive,
		&g.RedemptionLimit, &g.RedemptionCount, &st, &exp)
	if err != nil {
		c.errorLogger.Println(err)
		return nil
	}

	g.StartTime = st.Format("2006-01-02 15:04:05")
	g.ExpirationTime = exp.Format("2006-01-02 15:04:05")

	return g
}

func (c *DBConfig) UpdateGiftCodeCount(ctx context.Context, gift *models2.GiftCodeDAO) error {
	PingPQDB()
	ctxx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	_, err := c.dB.ExecContext(ctxx, queries2.UpdateGiftCodeRedemptionCount,
		&gift.ID)
	if err != nil {
		c.errorLogger.Println(err)
		return err
	}

	c.infoLogger.Println("Count Updated")
	return nil
}
