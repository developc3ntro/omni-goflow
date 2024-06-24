package resumes

import (
	"encoding/json"

	"github.com/developc3ntro/omni-goflow/assets"
	"github.com/developc3ntro/omni-goflow/envs"
	"github.com/developc3ntro/omni-goflow/flows"
	"github.com/developc3ntro/omni-goflow/flows/events"
	"github.com/developc3ntro/omni-goflow/utils"
	"github.com/nyaruka/gocommon/jsonx"
)

func init() {
	registerType(TypeWaitTimeout, readWaitTimeoutResume)
}

// TypeWaitTimeout is the type for resuming a session when a wait has timed out
const TypeWaitTimeout string = "wait_timeout"

// WaitTimeoutResume is used when a session is resumed because a wait has timed out
//
//	{
//	  "type": "wait_timeout",
//	  "contact": {
//	    "uuid": "9f7ede93-4b16-4692-80ad-b7dc54a1cd81",
//	    "name": "Bob",
//	    "created_on": "2018-01-01T12:00:00.000000Z",
//	    "language": "fra",
//	    "fields": {"gender": {"text": "Male"}},
//	    "groups": []
//	  },
//	  "resumed_on": "2000-01-01T00:00:00.000000000-00:00"
//	}
//
// @resume wait_timeout
type WaitTimeoutResume struct {
	baseResume
}

// NewWaitTimeout creates a new timeout resume with the passed in values
func NewWaitTimeout(env envs.Environment, contact *flows.Contact) *WaitTimeoutResume {
	return &WaitTimeoutResume{
		baseResume: newBaseResume(TypeWaitTimeout, env, contact),
	}
}

// Apply applies our state changes and saves any events to the run
func (r *WaitTimeoutResume) Apply(run flows.Run, logEvent flows.EventCallback) {
	logEvent(events.NewWaitTimedOut())

	r.baseResume.Apply(run, logEvent)
}

var _ flows.Resume = (*WaitTimeoutResume)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func readWaitTimeoutResume(sessionAssets flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Resume, error) {
	e := &baseResumeEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	r := &WaitTimeoutResume{}

	if err := r.unmarshal(sessionAssets, e, missing); err != nil {
		return nil, err
	}

	return r, nil
}

// MarshalJSON marshals this resume into JSON
func (r *WaitTimeoutResume) MarshalJSON() ([]byte, error) {
	e := &baseResumeEnvelope{}

	if err := r.marshal(e); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
