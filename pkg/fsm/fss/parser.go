package fss

import (
	"encoding/json"
	"fmt"
)

// parse fss language parser
var parse = fssParse

func (f *fssSymType) String() string {
	bs, _ := json.Marshal(f.Steps)
	return fmt.Sprintf("flow=%s,steps=%s", f.Flow, string(bs))
}

type StepType uint8

const (
	Normal StepType = iota
	Decision
)

type Step struct {
	Name     string `json:"name"`
	Action   `json:"action"`
	Returns  `json:"returns"`
	StepType `json:"step_type"`
}

type Returns []Return

type Return struct {
	State string `json:"state"`
	Next  string `json:"next"`
}

type ParamType uint8

const (
	StringType ParamType = iota
	NumberType
)

type Param struct {
	Name      string `json:"name"`
	ParamType `json:"type"`
	Value     interface{} `json:"value"`
}

type Action struct {
	Name string  `json:"name"`
	Args []Param `json:"args"`
}

type FlowStmt struct {
	*fssSymType
}

var flowSymTypePool = make(map[string]fssSymType)

func flowSymPoolPut(name string, sst fssSymType) {
	flowSymTypePool[name] = sst
}

func flowSymPoolGet(name string) (*FlowStmt, error) {
	sst, exist := flowSymTypePool[name]
	if !exist {
		return nil, fmt.Errorf("flow %s fssSymType not exist", name)
	}
	return &FlowStmt{&sst}, nil
}

type FlowRunStmt struct {
	*fssSymType
}

var flowRunSymTypePool = make(map[string]fssSymType)

func flowRunSymPoolPut(name string, sst fssSymType) {
	flowRunSymTypePool[name] = sst
}

func flowRunSymPoolGet(name string) (*FlowRunStmt, error) {
	sst, exist := flowRunSymTypePool[name]
	if !exist {
		return nil, fmt.Errorf("flowRun %s fssSymType not exist", name)
	}
	return &FlowRunStmt{&sst}, nil
}

type ActionMethodType uint8

const (
	ActionHTTPMethod ActionMethodType = iota
	ActionGRPCMethod
)

type ActionStatement struct {
	Name    string           `json:"name"`
	Addr    []string         `json:"addr"`
	Type    ActionMethodType `json:"type"`
	Args    []Param          `json:"args"`
	Returns Returns          `json:"returns"`
}

type ActionStmt struct {
	*fssSymType
}

var actionSymTypePool = make(map[string]fssSymType)

func actionSymPoolPut(name string, sst fssSymType) {
	actionSymTypePool[name] = sst
}

func actionSymPoolGet(name string) (*ActionStmt, error) {
	sst, exist := actionSymTypePool[name]
	if !exist {
		return nil, fmt.Errorf("action %s fssSymType not exist", name)
	}
	return &ActionStmt{&sst}, nil
}
