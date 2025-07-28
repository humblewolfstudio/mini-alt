package storage

import (
	"database/sql"
	"errors"
	"mini-alt/types"
	"mini-alt/utils"
	"time"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`PRAGMA foreign_keys = ON`)
	if err != nil {
		return nil, err
	}

	store := &SQLiteStore{db: db}
	if err := store.initSchema(); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *SQLiteStore) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS buckets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		created_at DATETIME NOT NULL
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

	CREATE TABLE IF NOT EXISTS object_metadata (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    object_id INTEGER NOT NULL,
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
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
`
	_, err := s.db.Exec(schema)
	return err
}

func (s *SQLiteStore) PutObject(bucket, object string, size int64) (Object, error) {
	now := time.Now()
	row := s.db.QueryRow(`
	INSERT INTO objects(bucket_name, key, size, last_modified)
	VALUES (?, ?, ?, ?)
	ON CONFLICT(bucket_name, key) DO UPDATE 
	SET size = excluded.size, last_modified = excluded.last_modified
	RETURNING id`,
		bucket, object, size, now)

	var id int64
	if err := row.Scan(&id); err != nil {
		return Object{}, err
	}

	return Object{Id: id, Key: object, Size: size, LastModified: now}, nil
}

func (s *SQLiteStore) PutBucket(bucket string) error {
	now := time.Now()
	_, err := s.db.Exec(`INSERT INTO buckets(name, created_at) VALUES (?, ?)`, bucket, now)
	if err != nil {
		return err
	}

	return nil
}

func (s *SQLiteStore) PutMetadata(objectId int64, metadata types.Metadata) error {
	_, err := s.db.Exec(`
		INSERT INTO object_metadata (
			object_id,
			acl,
			cache_control,
			content_disposition,
			content_encoding,
			content_language,
			content_length,
			content_md5,
			content_type,
			expires
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		objectId,
		metadata.ACL,
		metadata.CacheControl,
		metadata.ContentDisposition,
		metadata.ContentEncoding,
		metadata.ContentLanguage,
		metadata.ContentLength,
		metadata.ContentMD5,
		metadata.ContentType,
		metadata.Expires,
	)
	return err
}

func (s *SQLiteStore) ListObjects(bucket string) ([]Object, error) {
	rows, err := s.db.Query(`SELECT key, size, last_modified FROM objects WHERE bucket_name = ?`, bucket)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var objs []Object
	for rows.Next() {
		var o Object
		if err := rows.Scan(&o.Key, &o.Size, &o.LastModified); err != nil {
			return nil, err
		}
		objs = append(objs, o)
	}
	return objs, nil
}

func (s *SQLiteStore) ListBuckets() ([]Bucket, error) {
	rows, err := s.db.Query(`SELECT * FROM buckets`)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var buckets []Bucket
	for rows.Next() {
		var b Bucket
		if err := rows.Scan(&b.Id, &b.Name, &b.CreatedAt); err != nil {
			println(err.Error())
			return nil, err
		}
		buckets = append(buckets, b)
	}
	return buckets, nil
}

func (s *SQLiteStore) GetObject(bucket, key string) (Object, error) {
	row := s.db.QueryRow(`
		SELECT key, size, last_modified FROM objects WHERE bucket_name = ? AND key = ?`,
		bucket, key)

	var obj Object
	obj.Key = key
	if err := row.Scan(&obj.Key, &obj.Size, &obj.LastModified); err != nil {
		return Object{}, errors.New("object not found")
	}
	return obj, nil
}

func (s *SQLiteStore) GetBucket(bucket string) (Bucket, error) {
	row := s.db.QueryRow(`SELECT * FROM buckets WHERE name = ?`, bucket)
	var b Bucket
	b.Name = bucket
	if err := row.Scan(&b.Id, &b.Name, &b.CreatedAt); err != nil {
		return Bucket{}, err
	}

	return b, nil
}

func (s *SQLiteStore) DeleteObject(bucket, key string) error {
	prefix := key + "%"
	_, err := s.db.Exec(`DELETE FROM objects WHERE bucket_name = ? AND key LIKE ? ESCAPE '\'`, bucket, prefix)
	return err
}

func (s *SQLiteStore) DeleteBucket(bucket string) error {
	_, err := s.db.Exec(`DELETE FROM buckets WHERE name = ?`, bucket)
	return err
}

func (s *SQLiteStore) CreateCredentials() (accessKey, secretKey string, err error) {
	accessKey = utils.GenerateRandomKey(16)
	secretKey = utils.GenerateRandomKey(32)

	encryptedSecret, err := utils.Encrypt(secretKey)
	if err != nil {
		return "", "", err
	}

	_, err = s.db.Exec(`INSERT INTO credentials (access_key, secret_key_encrypted) VALUES (?, ?)`, accessKey, encryptedSecret)
	if err != nil {
		return "", "", err
	}

	return accessKey, secretKey, nil
}

func (s *SQLiteStore) GetSecretKey(accessKey string) (string, error) {
	row := s.db.QueryRow(`SELECT secret_key_encrypted FROM credentials WHERE access_key = ?`, accessKey)
	var encrypted string
	if err := row.Scan(&encrypted); err != nil {
		return "", err
	}

	return utils.Decrypt(encrypted)
}

func (s *SQLiteStore) ListCredentials() ([]Credentials, error) {
	rows, err := s.db.Query(`SELECT access_key, created_at FROM credentials`)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var creds []Credentials
	for rows.Next() {
		var c Credentials
		if err := rows.Scan(&c.AccessKey, &c.CreatedAt); err != nil {
			return nil, err
		}

		creds = append(creds, c)
	}

	return creds, nil
}
