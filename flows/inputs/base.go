package inputs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type readFunc func(session flows.Session, data json.RawMessage) (flows.Input, error)

var registeredTypes = map[string]readFunc{}

func registerType(name string, f readFunc) {
	registeredTypes[name] = f
}

type baseInput struct {
	uuid      flows.InputUUID
	channel   flows.Channel
	createdOn time.Time
}

func (i *baseInput) UUID() flows.InputUUID  { return i.uuid }
func (i *baseInput) Channel() flows.Channel { return i.channel }
func (i *baseInput) CreatedOn() time.Time   { return i.createdOn }

// Resolve resolves the given key when this input is referenced in an expression
func (i *baseInput) Resolve(env utils.Environment, key string) types.XValue {
	switch key {
	case "uuid":
		return types.NewXText(string(i.uuid))
	case "created_on":
		return types.NewXDateTime(i.createdOn)
	case "channel":
		return i.channel
	}

	return types.NewXResolveError(i, key)
}

type baseInputEnvelope struct {
	UUID      flows.InputUUID         `json:"uuid"`
	Channel   *flows.ChannelReference `json:"channel,omitempty" validate:"omitempty,dive"`
	CreatedOn time.Time               `json:"created_on" validate:"required"`
}

// ReadInput reads an input from the given typed envelope
func ReadInput(session flows.Session, envelope *utils.TypedEnvelope) (flows.Input, error) {
	f := registeredTypes[envelope.Type]
	if f == nil {
		return nil, fmt.Errorf("unknown input type: %s", envelope.Type)
	}
	input, err := f(session, envelope.Data)
	if err != nil {
		return nil, fmt.Errorf("unable to read input[type=%s]: %s", envelope.Type, err)
	}
	return input, nil
}
