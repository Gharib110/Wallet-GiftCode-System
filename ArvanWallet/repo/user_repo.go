package repo

import (
	"context"
	"github.com/Gharib110/ArvanWallet/models"
	"github.com/Gharib110/ArvanWallet/queries"
	"time"
)

type Users interface {
	CreateUser(context.Context, *models.UserDAO) error
	GetUserInfoByPhoneNumber(context.Context, *models.UserDAO) *models.User
	UpdateUserCharge(context.Context, *models.UserDAO) error
}

func (c *DBConfig) CreateUser(ctx context.Context, user *models.UserDAO) error {
	PingPQDB()
	ctxx, cancel := context.WithTimeout(ctx, 1*time.Minute) // Time should not be hardCoded
	defer cancel()

	res, err := c.dB.ExecContext(ctxx, queries.InsertIntoUsersTable,
		&user.Name, &user.PhoneNumber, &user.Charge)
	if err != nil {
		c.errorLogger.Println(err)
		return err
	}

	c.infoLogger.Println("User is Created !")
	c.infoLogger.Println(res.LastInsertId())
	return nil
}

func (c *DBConfig) GetUserInfoByPhoneNumber(ctx context.Context, user *models.UserDAO) *models.User {
	PingPQDB()
	ctxx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	data := c.dB.QueryRowContext(ctxx,
		queries.SelectUserByPhoneNumber, user.PhoneNumber)
	if data.Err() != nil {
		c.errorLogger.Println(data.Err())
		return nil
	}
	u := &models.User{
		ID:          0,
		Name:        "",
		PhoneNumber: "",
		Charge:      0,
	}
	err := data.Scan(&u.ID, &u.Name, &u.PhoneNumber, &u.Charge)
	if err != nil {
		c.errorLogger.Println(data.Err())
		return nil
	}

	return u
}

func (c *DBConfig) UpdateUserCharge(ctx context.Context, user *models.UserDAO) error {
	PingPQDB()
	ctxx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	_, err := c.dB.ExecContext(ctxx,
		queries.UpdateUserChargeByID, &user.Charge, &user.ID)
	if err != nil {
		c.errorLogger.Println(err)
		return err
	}

	c.infoLogger.Println("User Charge Updated")
	return nil
}
