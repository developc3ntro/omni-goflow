package flows_test

import (
	"testing"

	"github.com/developc3ntro/omni-goflow/assets"
	"github.com/developc3ntro/omni-goflow/assets/static"
	"github.com/developc3ntro/omni-goflow/envs"
	"github.com/developc3ntro/omni-goflow/excellent/types"
	"github.com/developc3ntro/omni-goflow/flows"
	"github.com/developc3ntro/omni-goflow/test"

	"github.com/stretchr/testify/assert"
)

func TestGlobals(t *testing.T) {
	ga1 := static.NewGlobal("org_name", "Org Name", "U-Report")
	ga2 := static.NewGlobal("access_token", "Access Token", "674372272")

	ga := flows.NewGlobalAssets([]assets.Global{ga1, ga2})

	g1 := ga.Get("org_name")

	assert.Equal(t, "Org Name", g1.Name())
	assert.Equal(t, ga1, g1.Asset())
	assert.Equal(t, assets.NewGlobalReference("org_name", "Org Name"), g1.Reference())

	env := envs.NewBuilder().Build()

	// check use in expressions
	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{
		"__default__":  types.NewXText("Org Name: U-Report\nAccess Token: 674372272"),
		"access_token": types.NewXText("674372272"),
		"org_name":     types.NewXText("U-Report"),
	}), flows.Context(env, ga))
}
