package jobs

import (
	"database/sql"
	"errors"
	"fmt"
	"mini-alt/storage/db"
)

func LoadInitialData(store *db.Store) {
	isFirst, err := isFirstStartup(store)
	if err != nil {
		fmt.Printf("Error checking startup state: %v\n", err)
		return
	}

	if !isFirst {
		return
	}

	err = LoadAllHomeDirectory(store)
	if err != nil {
		fmt.Printf("Error loading home directory: %v\n", err)
		return
	}

	setLoadedInitialData(store)
}

func isFirstStartup(store *db.Store) (bool, error) {
	value, err := db.GetConfig[bool](store, "LOADED_INITIAL_DATA")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}
		return false, fmt.Errorf("failed to check startup state: %w", err)
	}
	return !value, nil
}

func setLoadedInitialData(store *db.Store) {
	err := db.SetConfig(store, "LOADED_INITIAL_DATA", "true")
	if err != nil {
		fmt.Printf("Error setting initial data: %v\n", err)
	}
}
