package modifiers

import (
	"encoding/json"

	"github.com/developc3ntro/omni-goflow/assets"
	"github.com/developc3ntro/omni-goflow/envs"
	"github.com/developc3ntro/omni-goflow/flows"
	"github.com/developc3ntro/omni-goflow/flows/events"
	"github.com/developc3ntro/omni-goflow/utils"
)

func init() {
	registerType(TypeName, readNameModifier)
}

// TypeName is the type of our name modifier
const TypeName string = "name"

// NameModifier modifies the name of a contact
type NameModifier struct {
	baseModifier

	Name string `json:"name"`
}

// NewName creates a new name modifier
func NewName(name string) *NameModifier {
	return &NameModifier{
		baseModifier: newBaseModifier(TypeName),
		Name:         name,
	}
}

// Apply applies this modification to the given contact
func (m *NameModifier) Apply(env envs.Environment, assets flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) {
	if contact.Name() != m.Name {
		// truncate value if necessary
		name := utils.Truncate(m.Name, env.MaxValueLength())

		contact.SetName(name)
		log(events.NewContactNameChanged(name))
		ReevaluateGroups(env, assets, contact, log)
	}
}

var _ flows.Modifier = (*NameModifier)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func readNameModifier(assets flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Modifier, error) {
	m := &NameModifier{}
	return m, utils.UnmarshalAndValidate(data, m)
}
