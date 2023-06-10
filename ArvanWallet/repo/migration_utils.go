package repo

import (
	"github.com/Gharib110/ArvanWallet/queries"
)

type MigrationUtils interface {
	createTables()
	deleteTables()
}

func (c *DBConfig) deleteTables() {
	_, err := c.dB.Exec(queries.DropTransactionsTable)
	if err != nil {
		c.errorLogger.Panicln(err)
		return
	} else {
		c.infoLogger.Println("TransactionTable is dropped !")
	}

	_, err = c.dB.Exec(queries.DropUsersTable)
	if err != nil {
		c.errorLogger.Panicln(err)
		return
	} else {
		c.infoLogger.Println("UsersTable is dropped !")
	}
}

func (c *DBConfig) createTables() {
	_, err := c.dB.Exec(queries.CreateUsersTable)
	if err != nil {
		c.errorLogger.Panicln(err)
		return
	} else {
		c.infoLogger.Println("UsersTable is Created")
	}
	_, err = c.dB.Exec(queries.CreateTransactionsTable)
	if err != nil {
		c.errorLogger.Panicln(err)
		return
	} else {
		c.infoLogger.Println("TransactionsTable is Created")
	}
}
