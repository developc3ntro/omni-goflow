package events

import (
	"encoding/json"

	"github.com/developc3ntro/omni-goflow/envs"
	"github.com/developc3ntro/omni-goflow/flows"
	"github.com/nyaruka/gocommon/jsonx"
)

func init() {
	registerType(TypeEnvironmentRefreshed, func() flows.Event { return &EnvironmentRefreshedEvent{} })
}

// TypeEnvironmentRefreshed is the type of our environment changed event
const TypeEnvironmentRefreshed string = "environment_refreshed"

// EnvironmentRefreshedEvent events are sent by the caller to tell the engine to update the session environment.
//
//	{
//	  "type": "environment_refreshed",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "environment": {
//	    "date_format": "YYYY-MM-DD",
//	    "time_format": "hh:mm",
//	    "timezone": "Africa/Kigali",
//	    "allowed_languages": ["eng", "fra"]
//	  }
//	}
//
// @event environment_refreshed
type EnvironmentRefreshedEvent struct {
	BaseEvent

	Environment json.RawMessage `json:"environment"`
}

// NewEnvironmentRefreshed creates a new environment changed event
func NewEnvironmentRefreshed(env envs.Environment) *EnvironmentRefreshedEvent {
	marshalled, _ := jsonx.Marshal(env)
	return &EnvironmentRefreshedEvent{
		BaseEvent:   NewBaseEvent(TypeEnvironmentRefreshed),
		Environment: marshalled,
	}
}

var _ flows.Event = (*EnvironmentRefreshedEvent)(nil)
