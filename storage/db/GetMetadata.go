package db

import "mini-alt/models"

func (s *Store) GetMetadata(objectId int64) (models.ObjectMetadata, error) {
	row := s.db.QueryRow(`SELECT cache_control, content_disposition, content_encoding, content_language, content_length, content_md5, content_type, expires FROM object_metadata WHERE object_id = ?`, objectId)

	var metadata models.ObjectMetadata
	if err := row.Scan(&metadata.CacheControl, &metadata.ContentDisposition, &metadata.ContentEncoding, &metadata.ContentLanguage, &metadata.ContentLength, &metadata.ContentMD5, &metadata.ContentType, &metadata.Expires); err != nil {
		return models.ObjectMetadata{}, err
	}

	return metadata, nil
}
