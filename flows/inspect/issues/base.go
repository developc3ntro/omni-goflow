package issues

import (
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
)

type reportFunc func(flows.SessionAssets, flows.Flow, []flows.ExtractedReference, func(flows.Issue))

var RegisteredTypes = map[string]reportFunc{}

// registers a new type of issue
func registerType(name string, report reportFunc) {
	RegisteredTypes[name] = report
}

// base of all issue types
type baseIssue struct {
	Type_        string           `json:"type"`
	NodeUUID_    flows.NodeUUID   `json:"node_uuid"`
	ActionUUID_  flows.ActionUUID `json:"action_uuid,omitempty"`
	Language_    envs.Language    `json:"language,omitempty"`
	Description_ string           `json:"description"`
}

// creates a new base issue
func newBaseIssue(typeName string, nodeUUID flows.NodeUUID, actionUUID flows.ActionUUID, language envs.Language, description string) baseIssue {
	return baseIssue{
		Type_:        typeName,
		NodeUUID_:    nodeUUID,
		ActionUUID_:  actionUUID,
		Language_:    language,
		Description_: description,
	}
}

// Type returns the type of this issue
func (p *baseIssue) Type() string { return p.Type_ }

// NodeUUID returns the UUID of the node where issue is found
func (p *baseIssue) NodeUUID() flows.NodeUUID { return p.NodeUUID_ }

// ActionUUID returns the UUID of the action where issue is found
func (p *baseIssue) ActionUUID() flows.ActionUUID { return p.ActionUUID_ }

// Language returns the translation language if the issue was found in a translation
func (p *baseIssue) Language() envs.Language { return p.Language_ }

// Description returns the description of the issue
func (p *baseIssue) Description() string { return p.Description_ }

// Check returns all issues in the given flow
func Check(sa flows.SessionAssets, flow flows.Flow, refs []flows.ExtractedReference) []flows.Issue {
	issues := make([]flows.Issue, 0)
	report := func(i flows.Issue) {
		issues = append(issues, i)
	}

	for _, fn := range RegisteredTypes {
		fn(sa, flow, refs, report)
	}

	return issues
}