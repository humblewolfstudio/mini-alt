package encoding

import "encoding/xml"

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
