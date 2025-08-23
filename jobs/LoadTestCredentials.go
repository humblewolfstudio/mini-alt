package jobs

import (
	"mini-alt/storage/db"
	"os"
)

//goland:noinspection ALL
func LoadTestCredentials(store *db.Store) {
	if exists, err := store.ExistsTestUser(); err != nil || exists {
		return
	}

	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	newId, err := store.RegisterUser("test", "test", accessKey, "", true)
	if err != nil {
		println(err.Error())
		return
	}

	err = store.ForceCredentials(accessKey, secretKey, newId)
	if err != nil {
		println(err.Error())
		return
	}
}
