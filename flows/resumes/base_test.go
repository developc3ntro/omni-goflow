package resumes_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/resumes"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/test"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResumeTypes(t *testing.T) {
	assetsJSON, err := ioutil.ReadFile("testdata/_assets.json")
	require.NoError(t, err)

	typeNames := make([]string, 0)
	for typeName := range resumes.RegisteredTypes() {
		typeNames = append(typeNames, typeName)
	}

	sort.Strings(typeNames)

	for _, typeName := range typeNames {
		testResumeType(t, assetsJSON, typeName)
	}
}

func testResumeType(t *testing.T, assetsJSON json.RawMessage, typeName string) {
	testPath := fmt.Sprintf("testdata/%s.json", typeName)
	testFile, err := ioutil.ReadFile(testPath)
	require.NoError(t, err)

	tests := []struct {
		Description   string              `json:"description"`
		Wait          json.RawMessage     `json:"wait,omitempty"`
		Resume        json.RawMessage     `json:"resume"`
		ReadError     string              `json:"read_error,omitempty"`
		Events        json.RawMessage     `json:"events,omitempty"`
		RunStatus     flows.RunStatus     `json:"run_status,omitempty"`
		SessionStatus flows.SessionStatus `json:"session_status,omitempty"`
	}{}

	err = jsonx.Unmarshal(testFile, &tests)
	require.NoError(t, err)

	defer dates.SetNowSource(dates.DefaultNowSource)
	defer uuids.SetGenerator(uuids.DefaultGenerator)

	for i, tc := range tests {
		dates.SetNowSource(dates.NewFixedNowSource(time.Date(2018, 10, 18, 14, 20, 30, 123456, time.UTC)))
		uuids.SetGenerator(uuids.NewSeededGenerator(12345))

		testName := fmt.Sprintf("test '%s' for resume type '%s'", tc.Description, typeName)

		testAssetsJSON := assetsJSON
		if tc.Wait != nil {
			testAssetsJSON = test.JSONReplace(assetsJSON, []string{"flows", "[0]", "nodes", "[0]", "router", "wait"}, tc.Wait)
		}

		// create session assets
		sa, err := test.CreateSessionAssets(testAssetsJSON, "")
		require.NoError(t, err, "unable to create session assets in %s", testName)

		resume, err := resumes.ReadResume(sa, tc.Resume, assets.PanicOnMissing)

		if tc.ReadError != "" {
			rootErr := errors.Cause(err)
			assert.EqualError(t, rootErr, tc.ReadError, "read error mismatch in %s", testName)
			continue
		} else {
			assert.NoError(t, err, "unexpected read error in %s", testName)
		}

		flow, err := sa.Flows().Get(assets.FlowUUID("ed352c17-191e-4e75-b366-1b2c54bb32d8"))
		require.NoError(t, err)

		// start a waiting session
		env := envs.NewBuilder().Build()
		eng := engine.NewBuilder().Build()
		contact := flows.NewEmptyContact(sa, "Bob", envs.Language("eng"), nil)
		trigger := triggers.NewBuilder(env, flow.Reference(), contact).Manual().Build()
		session, _, err := eng.NewSession(sa, trigger)
		require.NoError(t, err)
		require.Equal(t, flows.SessionStatusWaiting, session.Status())

		// resume with our resume...
		sprint, err := session.Resume(resume)

		actual := tc
		actual.RunStatus = session.Runs()[0].Status()
		actual.SessionStatus = session.Status()

		// re-marshal the resume
		actual.Resume, err = jsonx.Marshal(resume)
		require.NoError(t, err)

		actual.Events, _ = jsonx.Marshal(sprint.Events())

		if !test.UpdateSnapshots {
			// check resume marshalled correctly
			test.AssertEqualJSON(t, tc.Resume, actual.Resume, "marshal mismatch in %s", testName)

			// check statuses
			assert.Equal(t, tc.RunStatus, actual.RunStatus, "run status mismatch in %s", testName)
			assert.Equal(t, tc.SessionStatus, actual.SessionStatus, "session status mismatch in %s", testName)

			// check events generated by resume
			test.AssertEqualJSON(t, tc.Events, actual.Events, "events mismatch in %s", testName)
		} else {
			tests[i] = actual
		}
	}

	if test.UpdateSnapshots {
		actualJSON, err := jsonx.MarshalPretty(tests)
		require.NoError(t, err)

		err = ioutil.WriteFile(testPath, actualJSON, 0666)
		require.NoError(t, err)
	}
}

func TestReadResume(t *testing.T) {
	env := envs.NewBuilder().Build()

	missingAssets := make([]assets.Reference, 0)
	missing := func(a assets.Reference, err error) { missingAssets = append(missingAssets, a) }

	sessionAssets, err := engine.NewSessionAssets(env, static.NewEmptySource(), nil)
	require.NoError(t, err)

	// error if no type field
	_, err = resumes.ReadResume(sessionAssets, []byte(`{"foo": "bar"}`), missing)
	assert.EqualError(t, err, "field 'type' is required")

	// error if we don't recognize action type
	_, err = resumes.ReadResume(sessionAssets, []byte(`{"type": "do_the_foo", "foo": "bar"}`), missing)
	assert.EqualError(t, err, "unknown type: 'do_the_foo'")
}

func TestResumeContext(t *testing.T) {
	env := envs.NewBuilder().Build()
	msg := flows.NewMsgIn("605e6309-343b-4cac-8309-e1de4cadd7b5", urns.URN("tel:1234567890"), nil, "Hello", nil)
	resume := resumes.NewMsg(env, nil, msg)

	assert.Equal(t, map[string]types.XValue{
		"type": types.NewXText("msg"),
	}, resume.Context(env))
}
