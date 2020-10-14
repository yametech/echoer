package common

const (
	DefaultNamespace  = "echoer"
	EventCollection   = "event"
	ActionCollection  = "action"
	Step              = "step"
	FlowCollection    = "flow"
	FlowRunCollection = "flowrun"
)

type ActionParamIdType = string

const (
	FlowId    ActionParamIdType = "flowId"
	StepName  ActionParamIdType = "stepName"
	AckStates ActionParamIdType = "ackStates"
	UUID      ActionParamIdType = "uuid"
)

type Common struct {
	FlowId   string `json:"flowId"`
	StepName string `json:"stepName"`
	AckState string `json:"ackState"`
	UUID     string `json:"uuid"`
}
