package utils

import (
	"math"
	"math/rand"
	"statForMarket/internal/model"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var eventTypes = []string{"login", "logout", "signup", "create", "insert", "update", "delete"}

func GenerateEventID() int {
	return rand.Intn(math.MaxInt)
}

func GenerateEvents(count int) []*model.Event {
	events := make([]*model.Event, count)

	for i := 0; i < count; i++ {
		event := &model.Event{
			EventID:   rand.Intn(math.MaxInt),
			EventType: eventTypes[rand.Intn(len(eventTypes))],
			UserID:    rand.Intn(math.MaxInt8),
			EventTime: randate().Format("2006-01-02T15:04:05"),
			Payload:   randStringBytesRmndr(rand.Intn(20)),
		}
		events[i] = event
	}

	return events

}

func randStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func randate() time.Time {
	min := time.Date(2022, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2023, 9, 30, 23, 59, 59, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}
