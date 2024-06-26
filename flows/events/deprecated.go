package events

import (
	"github.com/developc3ntro/omni-goflow/assets"
	"github.com/developc3ntro/omni-goflow/flows"
)

func init() {
	registerType(TypeClassifierCalled, func() flows.Event { return &ClassifierCalledEvent{} })
}

// TypeClassifierCalled is our type for the classification event
const TypeClassifierCalled string = "classifier_called"

// ClassifierCalledEvent events have been replaced by service_called.
type ClassifierCalledEvent struct {
	BaseEvent

	Classifier *assets.ClassifierReference `json:"classifier" validate:"required"`
	HTTPLogs   []*flows.HTTPLog            `json:"http_logs"`
}
