package crons

import (
	"github.com/robfig/cron/v3"
	"mini-alt/storage"
)

func SetupDeleteExpiredCredentials(store storage.Store) {
	c := cron.New()
	_, err := c.AddFunc("0 2 * * *", func() { store.DeleteExpiredCredentials() })
	if err != nil {
		println("Error running DeleteExpiredCredentials cron: ", err.Error())
		return
	}
	c.Start()
}
