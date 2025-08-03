package crons

import (
	"github.com/robfig/cron/v3"
	"mini-alt/storage/db"
)

func SetupDeleteExpiredUsers(store *db.Store) {
	store.DeleteExpiredUsers()
	c := cron.New()
	_, err := c.AddFunc("0 2 * * *", func() { store.DeleteExpiredUsers() })
	if err != nil {
		println("Error running DeleteExpiredUsers cron: ", err.Error())
		return
	}
	c.Start()
}
