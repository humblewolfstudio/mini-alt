package events

import (
	"mini-alt/events/types"
	"mini-alt/storage/db"
)

const EventVersion = "0.1"
const EventSource = "minialt:s3"
const S3SchemaVersion = "1.0"
const ConfigurationId = "Config"

func handleEventObject(db *db.Store, evt *types.Event) {
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

		if eventBucket.Name != evt.Records[0].S3.Bucket.Name {
			continue
		}

		Pool.Submit(WebhookJob{Event: evt, Url: event.Endpoint, Token: event.Token})
	}
}

func handleEventGlobal(db *db.Store, evt *types.Event) {
	events, err := db.ListEvents()
	if err != nil {
		println("Error handling event: ", err.Error())
		return
	}

	for _, event := range events {
		if !event.Global {
			continue
		}

		Pool.Submit(WebhookJob{Event: evt, Url: event.Endpoint, Token: event.Token})
	}
}
