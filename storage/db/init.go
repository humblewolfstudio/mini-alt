package db

import (
	"database/sql"
	"fmt"
)

func newSQLiteStore(path string) (*Store, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`PRAGMA foreign_keys = ON`)
	if err != nil {
		return nil, err
	}

	store := &Store{db: db}
	if err := store.initSchema(); err != nil {
		return nil, err
	}
	if err := store.seedInitialData(); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *Store) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS buckets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		owner INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS objects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		bucket_name TEXT NOT NULL,
		key TEXT NOT NULL,
		size INTEGER,
		last_modified DATETIME,
		FOREIGN KEY(bucket_name) REFERENCES buckets(name) ON DELETE CASCADE,
		UNIQUE(bucket_name, key)
	);

	CREATE INDEX IF NOT EXISTS idx_objects_bucket_name ON objects(bucket_name);

	CREATE TABLE IF NOT EXISTS object_metadata (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    object_id INTEGER NOT NULL UNIQUE,
	    acl TEXT,
	    cache_control TEXT,
	    content_disposition TEXT,
	    content_encoding TEXT,
	    content_language TEXT,
	    content_length INTEGER,
	    content_md5 TEXT,
	    content_type TEXT,
	    expires DATETIME,
	    FOREIGN KEY(object_id) REFERENCES objects(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS  credentials (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		access_key TEXT UNIQUE NOT NULL,
		secret_key_encrypted TEXT NOT NULL,
		user BOOLEAN NOT NULL,
		expires_at DATE DEFAULT NULL,
		status BOOLEAN NOT NULL DEFAULT TRUE,
		name TEXT,
		description TEXT,
		owner INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS users (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    username TEXT UNIQUE NOT NULL UNIQUE,
	    password TEXT NOT NULL,
	    token TEXT NOT NULL,
		access_key TEXT NOT NULL,
		admin BOOLEAN NOT NULL DEFAULT FALSE,
	    expires_at DATE DEFAULT NULL,
	    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE TABLE IF NOT EXISTS config (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    key TEXT UNIQUE NOT NULL,
	    value TEXT NOT NULL,
	    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS events (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    name TEXT UNIQUE NOT NULL,
	    description TEXT,
	    bucket_id INTEGER NOT NULL,
	    endpoint TEXT NOT NULL,
	    token TEXT,
	    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
`
	_, err := s.db.Exec(schema)
	return err
}

func (s *Store) seedInitialData() error {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to count users: %w", err)
	}

	if count > 0 {
		return nil
	}
	accessKey, _, err := s.PutCredentials("", "", "", true, 1)
	if err != nil {
		return err
	}

	id, err := s.RegisterUser("admin", "admin", accessKey, "", true)
	if err != nil {
		return fmt.Errorf("failed to insert default user: %w", err)
	}

	err = s.AddCredentialsOwner(accessKey, id)
	if err != nil {
		return fmt.Errorf("failed to insert default user: %w", err)
	}

	return nil
}
