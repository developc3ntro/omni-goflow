package actions

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

// TypeReply is the type for reply actions
const TypeReply string = "reply"

// ReplyAction can be used to reply to the current contact in a flow. The text field may contain templates.
//
// A `broadcast_created` event will be created with the evaluated text.
//
// ```
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "reply",
//     "text": "Hi @contact.name, are you ready to complete today's survey?",
//     "attachments": [],
//     "all_urns": false
//   }
// ```
//
// @action reply
type ReplyAction struct {
	BaseAction
	Text         string   `json:"text"`
	Attachments  []string `json:"attachments"`
	QuickReplies []string `json:"quick_replies,omitempty"`
	AllURNs      bool     `json:"all_urns,omitempty"`
}

type msgDestination struct {
	urn     urns.URN
	channel flows.Channel
}

// Type returns the type of this action
func (a *ReplyAction) Type() string { return TypeReply }

// Validate validates our action is valid and has all the assets it needs
func (a *ReplyAction) Validate(assets flows.SessionAssets) error {
	return nil
}

// Execute runs this action
func (a *ReplyAction) Execute(run flows.FlowRun, step flows.Step, log flows.EventLog) error {
	evaluatedText, evaluatedAttachments, evaluatedQuickReplies := a.evaluateMessage(run, step, a.Text, a.Attachments, a.QuickReplies, log)

	channelSet, err := run.Session().Assets().GetChannelSet()
	if err != nil {
		return err
	}

	preferredChannel := run.Contact().Channel()
	sendChannels := channelSet.WithRole(flows.ChannelRoleSend)

	// if contact has a preferred channel which can send, it takes priority
	if preferredChannel != nil && preferredChannel.HasRole(flows.ChannelRoleSend) {
		sendChannels = append([]flows.Channel{preferredChannel}, sendChannels...)
	}

	destinations := []msgDestination{}

	if a.AllURNs {
		// send to any URN which has a corresponding channel (i.e. is sendable)
		for _, u := range run.Contact().URNs() {
			channel := flows.GetChannelForURN(sendChannels, u)
			if channel != nil {
				destinations = append(destinations, msgDestination{urn: u.URN, channel: channel})
			}
		}
	} else {
		// send to first URN which has a corresponding channel (i.e. is sendable)
		for _, u := range run.Contact().URNs() {
			channel := flows.GetChannelForURN(sendChannels, u)
			if channel != nil {
				destinations = append(destinations, msgDestination{urn: u.URN, channel: channel})
				break
			}
		}
	}

	// create a new message for each URN+channel destination
	for _, dest := range destinations {
		msg := flows.NewMsgOut(dest.urn, dest.channel, evaluatedText, evaluatedAttachments, evaluatedQuickReplies)
		log.Add(events.NewMsgCreatedEvent(msg))
	}

	return nil
}
