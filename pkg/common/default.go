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
	FlowId          ActionParamIdType = "flowId"
	StepName        ActionParamIdType = "stepName"
	AckStates       ActionParamIdType = "ackStates"
	UUID            ActionParamIdType = "uuid"
	GlobalVariables ActionParamIdType = "globalVariables"
	CaKey           ActionParamIdType = "cakey"
	CaPEM           ActionParamIdType = "capem"
)

type Common struct {
	FlowId   string `json:"flowId"`
	StepName string `json:"stepName"`
	AckState string `json:"ackState"`
	UUID     string `json:"uuid"`
	// add data info and globalVariables
	GlobalVariables map[string]interface{} `json:"globalVariables"`
	Data            string                 `json:"data"`
}
