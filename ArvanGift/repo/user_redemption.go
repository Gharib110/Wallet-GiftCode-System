package repo

import (
	"container/list"
	"context"
	"database/sql"
	"github.com/Gharib110/ArvanGift/models"
	"github.com/Gharib110/ArvanGift/queries"
	"time"
)

type Redemption interface {
	CreateRedemption(context.Context, *models2.UserRedemptionDAO) error
	GetRedeemedRecords(context.Context) (*list.List, error)
}

// CreateRedemption create a new redemption have three types Created, Done, Failed
func (c *DBConfig) CreateRedemption(ctx context.Context, rd *models2.UserRedemptionDAO) error {
	PingPQDB()
	ctxx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	_, err := c.dB.ExecContext(ctxx, queries2.InsertIntoRedemptionTable,
		&rd.UserID, &rd.GiftCodeID, &rd.Type, &rd.RedeemedAt)
	if err != nil {
		c.errorLogger.Println(err)
		return err
	}

	c.infoLogger.Println("Redemption is created !")
	return nil
}

// GetRedeemedRecords Get the last 50 records of Redemptions
func (c *DBConfig) GetRedeemedRecords(ctx context.Context) (*list.List, error) {
	ctxx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	rows, err := c.dB.QueryContext(ctxx, queries2.GetRedeemedRecords)
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			c.errorLogger.Println(err)
			return
		}
	}(rows)

	if err != nil {
		c.errorLogger.Println(err)
		return nil, err
	}
	lst := list.New()

	for rows.Next() {
		data := models2.UserRedemption{
			ID:         0,
			UserID:     0,
			GiftCodeID: 0,
			RedeemedAt: time.Time{},
			Type:       "",
		}
		err = rows.Scan(&data.ID, &data.UserID,
			&data.GiftCodeID, &data.Type, &data.RedeemedAt)
		if err != nil {
			c.errorLogger.Println(err)
			return nil, err
		}
		lst.PushBack(&data)
	}

	c.infoLogger.Println(lst.Len())
	return lst, nil
}
