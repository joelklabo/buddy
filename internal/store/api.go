package store

import "time"

// StoreAPI defines the persistence operations used by the runner.
type StoreAPI interface {
	SaveActive(pubkey, sessionID string) error
	ClearActive(pubkey string) error
	Active(pubkey string) (SessionState, bool, error)

	LastCursor(pubkey string) (time.Time, error)
	SaveCursor(pubkey string, ts time.Time) error

	AlreadyProcessed(eventID string) (bool, error)
	MarkProcessed(eventID string) error

	RecentMessageSeen(pubkey, message string, window time.Duration) (bool, error)
}
