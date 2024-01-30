package postgres

import (
	"github.com/go-pg/pg/v10"
	"treasury/src/config"
)

func StartDB() (*pg.DB, error) {
	var (
		opts *pg.Options
		err  error
	)

	opts, err = pg.ParseURL(config.Cfg.DB.DSN)
	if err != nil {
		return nil, err
	}

	//connect db
	db := pg.Connect(opts)

	return db, err
}
