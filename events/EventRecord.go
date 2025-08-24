package events

import (
	"mini-alt/events/types"
	"mini-alt/models"
	"mini-alt/storage/db"
	"time"
)

func HandleEventBucket(db *db.Store, eventName types.EventName, bucket, accessKey, clientIp string) {
	bucketType := getBucket(bucket, accessKey)
	s3 := getS3(bucketType, nil, nil)
	reqParams := getRequestParameters(clientIp)
	record := getEventRecord(eventName, s3, reqParams, accessKey)

	records := make([]types.Record, 0)
	records = append(records, record)

	event := types.Event{
		Records: records,
	}

	handleEventGlobal(db, &event)
}

func HandleEventObject(db *db.Store, eventName types.EventName, bucket, key, etag string, size int64, accessKey, clientIp string) {
	bucketType := getBucket(bucket, accessKey)
	objectType := getObject(key, &etag, "", &size)
	s3 := getS3(bucketType, objectType, nil)
	reqParams := getRequestParameters(clientIp)
	record := getEventRecord(eventName, s3, reqParams, accessKey)

	records := make([]types.Record, 0)
	records = append(records, record)

	event := types.Event{
		Records: records,
	}

	handleEventObject(db, &event)
}

func HandleEventObjectDelete(db *db.Store, eventName types.EventName, bucket, key, accessKey, clientIp string) {
	bucketType := getBucket(bucket, accessKey)
	objectType := getObject(key, nil, "", nil)
	s3 := getS3(bucketType, objectType, nil)
	reqParams := getRequestParameters(clientIp)
	record := getEventRecord(eventName, s3, reqParams, accessKey)

	records := make([]types.Record, 0)
	records = append(records, record)

	event := types.Event{
		Records: records,
	}

	handleEventObject(db, &event)
}

func HandleEventObjectCopy(db *db.Store, eventName types.EventName, srcBucket, srcKey, dstBucket, dstKey, etag string, size int64, accessKey, clientIp string) {
	bucketType := getBucket(dstBucket, accessKey)
	objectType := getObject(dstKey, &etag, "", &size)
	s3 := getS3(bucketType, objectType, nil)

	source := getSource(srcBucket, srcKey)
	reqParams := getRequestParameters(clientIp)
	record := getEventRecord(eventName, s3, reqParams, accessKey)
	record.Source = source

	records := make([]types.Record, 0)
	records = append(records, record)

	event := types.Event{
		Records: records,
	}

	handleEventObject(db, &event)
}

func HandleEventObjectList(db *db.Store, eventName types.EventName, bucket, prefix string, maxKeys int64, content []models.Object, accessKey, clientIp string) {
	bucketType := getBucket(bucket, accessKey)
	objects := make([]types.Object, 0)
	for _, record := range content {
		objects = append(objects, types.Object{
			Key:       record.Key,
			ETag:      &record.ETag,
			Sequencer: "",
			Size:      &record.Size,
		})
	}

	s3 := getS3(bucketType, nil, &objects)
	reqParams := getRequestParameters(clientIp)
	record := getEventRecord(eventName, s3, reqParams, accessKey)
	record.RequestParameters.MaxKeys = &maxKeys
	record.RequestParameters.Prefix = &prefix

	records := make([]types.Record, 0)
	records = append(records, record)

	event := types.Event{
		Records: records,
	}

	handleEventObject(db, &event)
}

func getEventRecord(eventName types.EventName, s3 types.S3, requestParameters types.RequestParameters, accessKey string) types.Record {
	return types.Record{
		EventVersion:      EventVersion,
		EventSource:       EventSource,
		AwsRegion:         "us-east-1",
		EventTime:         time.Now().Format(time.RFC3339),
		EventName:         eventName,
		UserIdentity:      getUserIdentity(accessKey),
		RequestParameters: requestParameters,
		ResponseElements:  types.ResponseElements{},
		S3:                s3,
	}
}

func getS3(bucket *types.Bucket, object *types.Object, objects *[]types.Object) types.S3 {
	return types.S3{
		S3SchemaVersion: S3SchemaVersion,
		ConfigurationId: ConfigurationId,
		Bucket:          bucket,
		Object:          object,
		Objects:         objects,
	}
}

// GET THE OWNER OF THE BUCKET ACCESS KEY
func getBucket(bucket, accessKey string) *types.Bucket {
	return &types.Bucket{
		Name:          bucket,
		OwnerIdentity: getUserIdentity(accessKey),
		Arn:           "", // TODO add arn
	}
}

func getSource(srcBucket, srcKey string) *types.Source {
	return &types.Source{
		Bucket: srcBucket,
		Key:    srcKey,
	}
}

// TODO add sequencer???
func getObject(key string, eTag *string, sequencer string, size *int64) *types.Object {
	return &types.Object{
		Key:       key,
		ETag:      eTag,
		Size:      size,
		Sequencer: sequencer,
	}
}

func getUserIdentity(accessKey string) types.UserIdentity {
	return types.UserIdentity{
		PrincipalId: accessKey,
	}
}

func getRequestParameters(clientIp string) types.RequestParameters {
	return types.RequestParameters{
		SourceIpAddress: clientIp,
	}
}
