package factory

import (
	"encoding/json"
	"testing"

	"github.com/yametech/echoer/pkg/core"
	"github.com/yametech/echoer/pkg/fsm/fss"
	"github.com/yametech/echoer/pkg/resource"
)

var _ IStore = &FakeStoreImpl{}

type FakeStoreImpl struct {
	data string
}

func (f *FakeStoreImpl) GetFlowRun(s2 string) (*resource.FlowRun, error) {
	panic("implement me")
}

func (f *FakeStoreImpl) CreateFlow(fr *resource.Flow) error {
	panic("implement me")
}

func (f *FakeStoreImpl) CreateAction(fr *resource.Action) error {
	panic("implement me")
}

func (f *FakeStoreImpl) GetAction(s2 string) (*resource.Action, error) {
	return &resource.Action{
		Metadata: core.Metadata{
			Name: s2,
			Kind: resource.ActionKind,
		},
		Spec: resource.ActionSpec{
			System:    "compass",
			ServeType: resource.HTTP,
			Endpoints: []string{"http://127.0.0.1:8081"},
			Params: map[resource.ParamNameType]resource.ParamType{
				resource.ParamNameType("pipeline"): resource.STR,
			},
			ReturnStates: []string{"SUCCESS", "FAILED"},
		},
	}, nil
}

func (f *FakeStoreImpl) CreateFlowRun(fr *resource.FlowRun) error {
	bs, _ := json.Marshal(fr)
	f.data = string(bs)
	return nil
}

func NewFakeStoreImpl() *FakeStoreImpl {
	return &FakeStoreImpl{}
}

const fsl = `
flow_run test_fsl_parse
	step A => (SUCCESS->stopped | FAILED-> A){
		action = "xx";
		args = (pipeline="abc");
	};
flow_run_end
`

func TestTranslation_ToFlowRun(t *testing.T) {
	store := NewFakeStoreImpl()
	_ = NewTranslation(store)
	_, err := fss.NewFlowRunFSLParser().Parse(fsl)
	if err != nil {
		t.Fatal(err)
	}

	_ = store.data
}
