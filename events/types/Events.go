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
	EventName EventName
	AccessKey string
}

type Eventer interface {
	GetBase() Event
}
