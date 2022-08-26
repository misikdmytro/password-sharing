package repository

import (
	"database/sql"

	"github.com/misikdmitriy/password-sharing/models"

	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
)

func InitDB(connection string) (*gorp.DbMap, error) {
	db, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, err
	}

	dbmap := &gorp.DbMap{
		Db:      db,
		Dialect: gorp.PostgresDialect{},
	}

	table := dbmap.AddTableWithName(models.Password{}, "tbl_passwords")
	table.ColMap("Id").Rename("id")

	return dbmap, nil
}
