package config

import (
	"log"
	"os"
)

var Cfg *config

func init() {
	log.Println("Config init...")
	Cfg = &config{}
	if err := Cfg.Load(); err != nil {
		log.Fatalln(err)
	}
}

type config struct {
	DB struct {
		DSN string `env:"DB_DSN,required"`
	}
}

func (c *config) Load() error {
	dsn, exists := os.LookupEnv("DB_DSN")

	if !exists {
		// Print the value of the environment variable
		log.Fatalln("Don't set env DB_DSN")
	}
	c.DB.DSN = dsn

	return nil
}
