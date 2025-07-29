package crons

import "mini-alt/storage"

func StartupCronJobs(store storage.Store) {
	SetupDeleteExpiredCredentials(store)
	SetupDeleteExpiredUsers(store)
}
