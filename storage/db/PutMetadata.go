package db

import "mini-alt/types"

func (s *Store) PutMetadata(objectId int64, metadata types.Metadata) error {
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
