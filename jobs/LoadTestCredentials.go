package jobs

import (
	"mini-alt/storage/db"
)

//goland:noinspection ALL
func LoadTestCredentials(store *db.Store) {
	if exists, err := store.ExistsTestUser(); err != nil || exists {
		return
	}

	accessKey, secretKey := "zAYWfUfum6mKhLbK", "mNcRRjk4Uy9eFnzIUHANbnRtnKQsFi2I"

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
