package model

type Event struct {
	EventID   int    `json:"eventID"`
	EventType string `json:"eventType"`
	UserID    int    `json:"userID"`
	EventTime string `json:"eventTime"`
	Payload   string `json:"payload"`
}

type EventFilter struct {
	EventType string `json:"eventType"`
	From      string `json:"from"`
	To        string `json:"to"`
}
