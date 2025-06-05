package encoding

import (
	"encoding/xml"
	"time"
)

type BucketXML struct {
	XMLName   xml.Name `xml:"Bucket"`
	Name      string   `xml:"Name"`
	CreatedAt string   `xml:"CreationDate"`
}

type ListAllMyBucketsResult struct {
	XMLName xml.Name `xml:"ListAllMyBucketsResult"`
	XMLNS   string   `xml:"xmlns,attr"`
	Buckets struct {
		Bucket []BucketXML `xml:"Bucket"`
	} `xml:"Buckets"`
}

type S3Error struct {
	XMLName    xml.Name `xml:"Error"`
	Code       string   `xml:"Code"`
	Message    string   `xml:"Message"`
	BucketName string   `xml:"BucketName,omitempty"`
	RequestID  string   `xml:"RequestId"`
	HostID     string   `xml:"HostId"`
}

type ListBucketResult struct {
	XMLName     xml.Name  `xml:"ListBucketResult"`
	Contents    []Content `xml:"Contents"`
	IsTruncated bool      `xml:"IsTruncated"`
	KeyCount    int       `xml:"KeyCount"`
	MaxKeys     int       `xml:"MaxKeys"`
	Name        string    `xml:"Name"`
	StartAfter  string    `xml:"StartAfter"`
}

type Content struct {
	Key          string    `xml:"Key"`
	LastModified time.Time `xml:"LastModified"`
	Size         int64     `xml:"Size"`
}
