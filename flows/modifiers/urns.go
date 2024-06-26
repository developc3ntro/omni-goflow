package modifiers

import (
	"encoding/json"

	"github.com/developc3ntro/omni-goflow/assets"
	"github.com/developc3ntro/omni-goflow/envs"
	"github.com/developc3ntro/omni-goflow/flows"
	"github.com/developc3ntro/omni-goflow/flows/events"
	"github.com/developc3ntro/omni-goflow/utils"
	"github.com/nyaruka/gocommon/urns"
)

func init() {
	registerType(TypeURNs, readURNsModifier)
}

// TypeURNs is the type of our URNs modifier
const TypeURNs string = "urns"

// URNsModification is the type of modification to make
type URNsModification string

// the supported types of modification
const (
	URNsAppend URNsModification = "append"
	URNsRemove URNsModification = "remove"
	URNsSet    URNsModification = "set"
)

// URNsModifier modifies the URNs on a contact
type URNsModifier struct {
	baseModifier

	URNs         []urns.URN       `json:"urns" validate:"required"`
	Modification URNsModification `json:"modification" validate:"required,eq=append|eq=remove|eq=set"`
}

// NewURNs creates a new URNs modifier
func NewURNs(urnz []urns.URN, modification URNsModification) *URNsModifier {
	return &URNsModifier{
		baseModifier: newBaseModifier(TypeURNs),
		URNs:         urnz,
		Modification: modification,
	}
}

// Apply applies this modification to the given contact
func (m *URNsModifier) Apply(env envs.Environment, assets flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) {
	modified := false

	if m.Modification == URNsSet {
		modified = contact.ClearURNs()
	}

	for _, urn := range m.URNs {
		// normalize the URN
		urn := urn.Normalize(string(env.DefaultCountry()))

		if err := urn.Validate(); err != nil {
			log(events.NewErrorf("'%s' is not valid URN", urn))
		} else {
			if m.Modification == URNsAppend || m.Modification == URNsSet {
				modified = contact.AddURN(urn, nil)
			} else {
				modified = contact.RemoveURN(urn)
			}
		}
	}

	if modified {
		log(events.NewContactURNsChanged(contact.URNs().RawURNs()))
		ReevaluateGroups(env, assets, contact, log)
	}
}

var _ flows.Modifier = (*URNsModifier)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func readURNsModifier(assets flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Modifier, error) {
	m := &URNsModifier{}
	return m, utils.UnmarshalAndValidate(data, m)
}
