package db

import (
	"mini-alt/models"
	"mini-alt/types"
)

func (s *Store) getMetadata(objectId int64) (models.ObjectMetadata, error) {
	row := s.db.QueryRow(`SELECT cache_control, content_disposition, content_encoding, content_language, content_length, content_md5, content_type, expires FROM object_metadata WHERE object_id = ?`, objectId)

	var metadata models.ObjectMetadata
	if err := row.Scan(&metadata.CacheControl, &metadata.ContentDisposition, &metadata.ContentEncoding, &metadata.ContentLanguage, &metadata.ContentLength, &metadata.ContentMD5, &metadata.ContentType, &metadata.Expires); err != nil {
		return models.ObjectMetadata{}, err
	}

	return metadata, nil
}

func (s *Store) putMetadata(objectId int64, metadata types.Metadata) error {
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
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT (object_id) DO UPDATE SET
		acl = excluded.acl,
		cache_control = excluded.cache_control,
		content_disposition = excluded.content_disposition,
		content_encoding = excluded.content_encoding,
		content_language = excluded.content_language,
		content_length = excluded.content_length,
		content_md5 = excluded.content_md5,
		content_type = excluded.content_type,
		expires = excluded.expires`,
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

func (s *Store) copyMetadata(oldObjectId, newObjectId int64) error {
	_, err := s.db.Exec(`
	INSERT OR REPLACE INTO object_metadata (
		object_id, acl, cache_control, content_disposition, content_encoding, 
		content_language, content_length, content_md5, content_type, expires
	)
	SELECT 
		?, acl, cache_control, content_disposition, content_encoding, 
		content_language, content_length, content_md5, content_type, expires 
	FROM object_metadata 
	WHERE object_id = ?
`, newObjectId, oldObjectId)

	return err
}
