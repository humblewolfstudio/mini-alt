package crons

import (
	"mini-alt/storage/db"
)

func StartupCronJobs(store *db.Store) {
	SetupDeleteExpiredCredentials(store)
	SetupDeleteExpiredUsers(store)
}
