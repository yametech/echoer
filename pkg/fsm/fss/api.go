package fss

import (
	"fmt"
	"sync"

	"github.com/yametech/echoer/pkg/resource"
)

type FlowRunFSLParser struct {
	sync.Mutex
	flowRun *resource.FlowRun
	symPool map[string]fssSymType
}

func NewFlowRunFSLParser() *FlowRunFSLParser {
	return &FlowRunFSLParser{
		Mutex:   sync.Mutex{},
		symPool: flowRunSymTypePool,
	}
}

func (f *FlowRunFSLParser) Parse(fsl string) (*FlowRunStmt, error) {
	defer f.Unlock()
	f.Lock()
	val := parse(NewFssLexer([]byte(fsl)))
	if val != 0 {
		return nil, fmt.Errorf("parse flow run (%s) error", fsl)
	}
	var sst fssSymType
	for k, v := range f.symPool {
		sst = v
		delete(f.symPool, k)
		break
	}
	return &FlowRunStmt{&sst}, nil
}

type FlowFSLParser struct {
	sync.Mutex
	flow    *resource.Flow
	symPool map[string]fssSymType
}

func NewFlowFSLParser() *FlowFSLParser {
	return &FlowFSLParser{
		Mutex:   sync.Mutex{},
		symPool: flowSymTypePool,
	}
}

func (f *FlowFSLParser) Parse(fsl string) (*FlowStmt, error) {
	defer f.Unlock()
	f.Lock()
	val := parse(NewFssLexer([]byte(fsl)))
	if val != 0 {
		return nil, fmt.Errorf("parse flow (%s) error", fsl)
	}
	var sst fssSymType
	for k, v := range f.symPool {
		sst = v
		delete(f.symPool, k)
		break
	}
	return &FlowStmt{&sst}, nil
}

type ActionFSLParser struct {
	sync.Mutex
	flow    *resource.Action
	symPool map[string]fssSymType
}

func NewActionFSLParser() *ActionFSLParser {
	return &ActionFSLParser{
		Mutex:   sync.Mutex{},
		symPool: actionSymTypePool,
	}
}

func (f *ActionFSLParser) Parse(fsl string) (*ActionStmt, error) {
	defer f.Unlock()
	f.Lock()
	val := parse(NewFssLexer([]byte(fsl)))
	if val != 0 {
		return nil, fmt.Errorf("parse flow (%s) error", fsl)
	}
	var sst fssSymType
	for k, v := range f.symPool {
		sst = v
		delete(f.symPool, k)
		break
	}
	return &ActionStmt{&sst}, nil
}
