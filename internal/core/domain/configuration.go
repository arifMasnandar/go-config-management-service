package domain

import "time"

// Config represents data about a record Config.
type Config struct {
	Name              string                 `json:"name"`
	Value             string                 `json:"value"`
	JsonValue         map[string]interface{} `json:"json_value"`
	Version           int                    `json:"version"`
	RollbackedVersion int                    `json:"rollbacked_version,omitempty"` // Optional field for copied version
	CreatedAt         time.Time              `json:"created_at,omitempty"`         // Optional field for creation timestamp
}
