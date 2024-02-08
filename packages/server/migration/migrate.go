package main

import (
	"area-server/db/postgres"
)

func main() {
	/* err := godotenv.Load("../../.env")
	if err != nil {
		panic(err)
	} */
	// Do nothing
	pg := postgres.Init()
	pg.Connect()

	pg.Migrate() // Migrate database - Change to false in config if you don't want to migrate
}
