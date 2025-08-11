package disk

import (
	"crypto/md5"
	"encoding/base64"
	"io"
	"io/fs"
	"mime"
	"mini-alt/storage/db"
	"mini-alt/types"
	"net/http"
	"os"
	"path/filepath"
)

// LoadAllHomeDirectory loads all existing buckets and their objects from disk into the database.
func LoadAllHomeDirectory(s *db.Store) error {
	bucketsDir, err := GetBucketsDir()
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(bucketsDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		bucketName := entry.Name()

		err := createBucket(s, bucketName)
		if err != nil {
			return err
		}

		bucketPath := filepath.Join(bucketsDir, bucketName)

		err = filepath.WalkDir(bucketPath, func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}

			if path == bucketPath {
				return nil
			}

			objectKey, err := filepath.Rel(bucketPath, path)
			if err != nil {
				return err
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer func(file *os.File) {
				_ = file.Close()
			}(file)

			return createFile(s, bucketName, objectKey, path)
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func createBucket(s *db.Store, bucketName string) error {
	if err := s.PutBucket(bucketName, 1); err != nil {
		return err
	}

	if err := CreateBucket(bucketName); err != nil {
		return err
	}

	return nil
}

func createFile(s *db.Store, bucketName string, objectKey string, filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		object, err := s.PutObject(bucketName, objectKey, 0)
		if err != nil {
			return err
		}

		meta := types.Metadata{
			ContentLength: 0,
			ContentType:   "application/x-directory",
		}

		return s.PutMetadata(object.Id, meta)
	}

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	size := info.Size()

	buffer := make([]byte, 512)
	n, _ := f.Read(buffer)
	_, _ = f.Seek(0, io.SeekStart)

	contentType := http.DetectContentType(buffer[:n])
	if contentType == "application/octet-stream" {
		ext := filepath.Ext(objectKey)
		if mimeByExt := mime.TypeByExtension(ext); mimeByExt != "" {
			contentType = mimeByExt
		}
	}

	hasher := md5.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return err
	}
	md5sum := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	object, err := s.PutObject(bucketName, objectKey, size)
	if err != nil {
		return err
	}

	meta := types.Metadata{
		ContentLength: size,
		ContentType:   contentType,
		ContentMD5:    md5sum,
	}

	return s.PutMetadata(object.Id, meta)
}
