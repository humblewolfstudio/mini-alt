package types

type EventName string

//goland:noinspection GoUnusedConst
const (
	EventGet           EventName = "s3:ObjectAccessed:Get"
	EventHead          EventName = "s3:ObjectAccessed:Head"
	EventPut           EventName = "s3:ObjectCreated:Put"
	EventCopy          EventName = "s3:ObjectCreated:Copy"
	EventCopied        EventName = "s3:ObjectCreated:Copied"
	EventDelete        EventName = "s3:ObjectRemoved:Delete"
	EventGetPrefix     EventName = "s3:ObjectAccessed:*"
	EventDeletedPrefix EventName = "s3:ObjectRemoved:*"
	EventBucketCreated EventName = "s3:BucketCreated"
	EventBucketDeleted EventName = "s3:BucketRemoved"
)

type Event struct {
	Records []Record `json:"Records"`
}

type Record struct {
	EventVersion      string            `json:"eventVersion"`
	EventSource       string            `json:"eventSource"`
	AwsRegion         string            `json:"awsRegion"`
	EventTime         string            `json:"eventTime"`
	EventName         EventName         `json:"eventName"`
	UserIdentity      UserIdentity      `json:"userIdentity"`
	RequestParameters RequestParameters `json:"requestParameters"`
	ResponseElements  ResponseElements  `json:"responseElements"`
	S3                S3                `json:"s3"`
	Source            *Source           `json:"source,omitempty"`
}

type UserIdentity struct {
	PrincipalId string `json:"principalId"`
}

type RequestParameters struct {
	SourceIpAddress string  `json:"sourceIpAddress"`
	Prefix          *string `json:"prefix,omitempty"`
	MaxKeys         *int64  `json:"maxKeys,omitempty"`
}

type ResponseElements struct {
	XAmzRequestId string `json:"x-amz-request-id"`
	XAmzId2       string `json:"x-amz-id-2"`
}

type S3 struct {
	S3SchemaVersion string    `json:"s3SchemaVersion"`
	ConfigurationId string    `json:"configurationId"`
	Bucket          *Bucket   `json:"bucket,omitempty"`
	Object          *Object   `json:"object,omitempty"`
	Objects         *[]Object `json:"objects,omitempty"`
}

type Source struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}

type Bucket struct {
	Name          string       `json:"name"`
	OwnerIdentity UserIdentity `json:"ownerIdentity"`
	Arn           string       `json:"arn"`
}

type Object struct {
	Key       string  `json:"key"`
	Size      *int64  `json:"size,omitempty"`
	ETag      *string `json:"eTag,omitempty"`
	Sequencer string  `json:"sequencer"`
}

type Eventer interface {
	GetBase() Event
}
