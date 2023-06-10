package repo

import (
	"github.com/Gharib110/ArvanGift/queries"
)

type MigrationUtils interface {
	createTables()
	deleteTables()
}

func (c *DBConfig) deleteTables() {
	_, err := c.dB.Exec(queries2.DropUserRedemptionTable)
	if err != nil {
		c.errorLogger.Panicln(err)
		return
	} else {
		c.infoLogger.Println("UserRedemptionTable is dropped !")
	}

	_, err = c.dB.Exec(queries2.DropGiftCodeTable)
	if err != nil {
		c.errorLogger.Panicln(err)
		return
	} else {
		c.infoLogger.Println("GiftCodeTable is dropped !")
	}
}

func (c *DBConfig) createTables() {
	_, err := c.dB.Exec(queries2.CreateGiftCodesTable)
	if err != nil {
		c.errorLogger.Panicln(err)
		return
	} else {
		c.infoLogger.Println("GiftCodeTableTable is Created")
	}
	_, err = c.dB.Exec(queries2.CreateRedemptionTable)
	if err != nil {
		c.errorLogger.Panicln(err)
		return
	} else {
		c.infoLogger.Println("RedemptionTableTable is Created")
	}
}
