package main

import (
	"fmt"
	"github.com/go-pg/migrations/v8"
	"log"
	"treasury/src/domain"
)

func init() {
	log.Println("1 initial...")
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("creating table entries...")
		_, err := db.Exec(`CREATE TABLE IF NOT EXISTS entries(
			id serial PRIMARY KEY,
			uid varchar(45) NOT NULL,
		    last_name varchar(255),
		    first_name varchar(255))`)
		if err != nil {
			return err
		}
		fmt.Println("creating table states...")
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS states (
    		id serial PRIMARY KEY,
		  	info varchar(45) NOT NULL)`)
		if err != nil {
			return err
		}

		fmt.Println("set empty state...")
		state := new(domain.State)
		state.Info = "empty"
		_, err = db.Model(state).Insert()

		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping table entries...")
		_, err := db.Exec(`DROP TABLE entries`)
		if err != nil {
			return err
		}

		fmt.Println("dropping table states...")
		_, err = db.Exec(`DROP TABLE states`)

		return err
	})

}
