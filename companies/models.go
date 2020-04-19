package companies

import (
	"time"
)

type Status struct {
	ID                 string
	SenderCountryISO   string
	ReceiverCountryISO string
	History            []HistoryItem
}

type HistoryItem struct {
	At                 time.Time
	Location           string
	LocationCountryISO string
	Message            string
}
