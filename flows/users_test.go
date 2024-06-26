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

func TestUsers(t *testing.T) {
	ua1 := static.NewUser("bob@nyaruka.com", "Bob McTickets")
	ua2 := static.NewUser("jim@nyaruka.com", "")

	ua := flows.NewUserAssets([]assets.User{ua1, ua2})

	u1 := ua.Get("bob@nyaruka.com")

	assert.Equal(t, "Bob McTickets", u1.Format())
	assert.Equal(t, "Bob McTickets", u1.Name())
	assert.Equal(t, ua1, u1.Asset())
	assert.Equal(t, assets.NewUserReference("bob@nyaruka.com", "Bob McTickets"), u1.Reference())

	// nil object returns nil reference
	assert.Nil(t, (*flows.User)(nil).Reference())

	env := envs.NewBuilder().Build()

	// check use in expressions
	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{
		"__default__": types.NewXText("Bob McTickets"),
		"email":       types.NewXText("bob@nyaruka.com"),
		"name":        types.NewXText("Bob McTickets"),
		"first_name":  types.NewXText("Bob"),
	}), flows.Context(env, u1))

	u2 := ua.Get("jim@nyaruka.com")

	assert.Equal(t, "jim@nyaruka.com", u2.Format()) // fallsback on email

	// check use in expressions
	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{
		"__default__": types.NewXText("jim@nyaruka.com"),
		"email":       types.NewXText("jim@nyaruka.com"),
		"name":        types.NewXText(""),
		"first_name":  types.NewXText(""),
	}), flows.Context(env, u2))
}
