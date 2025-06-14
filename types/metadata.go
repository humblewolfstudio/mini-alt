package types

type Metadata struct {
	Id                 int64
	ObjectId           int64
	ACL                string
	CacheControl       string
	ContentDisposition string
	ContentEncoding    string
	ContentLanguage    string
	ContentLength      int64
	ContentMD5         string
	ContentType        string
	Expires            string // RFC1123 or RFC3339 format
}
