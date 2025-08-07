package db

func (s *Store) MetadataCopy(oldObjectId, newObjectId int64) error {
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
