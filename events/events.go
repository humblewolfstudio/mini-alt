package events

import (
	"mini-alt/events/types"
	"mini-alt/storage/db"
)

func HandleEventObject(db *db.Store, eventName types.EventName, key, bucket string, accessKey string) {
	events, err := db.ListEvents()
	if err != nil {
		println("Error handling event: ", err.Error())
		return
	}

	for _, event := range events {

		eventBucket, err := db.GetBucketById(event.BucketId)
		if err != nil {
			println("Error handling event: ", err.Error())
			return
		}

		if eventBucket.Name != bucket {
			continue
		}

		evt := types.ObjectEvent{
			Event: types.Event{
				EventName: eventName,
				AccessKey: accessKey,
			},
			Key: key,
		}

		Pool.Submit(WebhookJob{Event: evt, Link: event.Endpoint, Token: event.Token})
	}
}

func HandleEventBucket(db *db.Store, eventName types.EventName, bucket string, accessKey string) {
	events, err := db.ListEvents()
	if err != nil {
		println("Error handling event: ", err.Error())
		return
	}

	for _, event := range events {
		if !event.Global {
			continue
		}

		evt := types.GlobalEvent{
			Event: types.Event{
				EventName: eventName,
				AccessKey: accessKey,
			},
			Bucket: bucket,
		}

		Pool.Submit(WebhookJob{Event: evt, Link: event.Endpoint, Token: event.Token})
	}
}
