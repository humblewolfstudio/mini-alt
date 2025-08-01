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
