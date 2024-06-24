package test

import (
	"os"

	"github.com/developc3ntro/omni-goflow/assets"
	"github.com/developc3ntro/omni-goflow/assets/static"
	"github.com/developc3ntro/omni-goflow/envs"
	"github.com/developc3ntro/omni-goflow/flows"
	"github.com/developc3ntro/omni-goflow/flows/definition/migrations"
	"github.com/developc3ntro/omni-goflow/flows/engine"
	"github.com/nyaruka/gocommon/uuids"
)

// LoadSessionAssets loads a session assets instance from a static JSON file
func LoadSessionAssets(env envs.Environment, path string) (flows.SessionAssets, error) {
	assetsJSON, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	source, err := static.NewSource(assetsJSON)
	if err != nil {
		return nil, err
	}

	mconfig := &migrations.Config{BaseMediaURL: "http://temba.io/"}

	return engine.NewSessionAssets(env, source, mconfig)
}

func LoadFlowFromAssets(env envs.Environment, path string, uuid assets.FlowUUID) (flows.Flow, error) {
	sa, err := LoadSessionAssets(env, path)
	if err != nil {
		return nil, err
	}

	return sa.Flows().Get(uuid)
}

func NewChannel(name string, address string, schemes []string, roles []assets.ChannelRole, parent *assets.ChannelReference) *flows.Channel {
	return flows.NewChannel(static.NewChannel(assets.ChannelUUID(uuids.New()), name, address, schemes, roles, parent))
}

func NewTelChannel(name string, address string, roles []assets.ChannelRole, parent *assets.ChannelReference, country envs.Country, matchPrefixes []string, allowInternational bool) *flows.Channel {
	return flows.NewChannel(static.NewTelChannel(assets.ChannelUUID(uuids.New()), name, address, roles, parent, country, matchPrefixes, allowInternational))
}

func NewClassifier(name, type_ string, intents []string) *flows.Classifier {
	return flows.NewClassifier(static.NewClassifier(assets.ClassifierUUID(uuids.New()), name, type_, intents))
}
