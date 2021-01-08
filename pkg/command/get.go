package command

import (
	"encoding/json"
	"fmt"
	"github.com/yametech/echoer/pkg/common"
	"github.com/yametech/echoer/pkg/core"
	"github.com/yametech/echoer/pkg/resource"
	"github.com/yametech/echoer/pkg/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"strings"
)

type Get struct {
	storage.IStorage
}

func (g *Get) Name() string {
	return `GET`
}

func (g *Get) Execute(args ...string) Reply {
	if reply := checkArgsExpected(args, 2); reply != nil {
		return reply
	}
	resType := args[0]
	if storage.GetResourceCoder(resType) == nil {
		return &ErrorReply{Message: fmt.Sprintf("this type (%s) is not supported", resType)}
	}

	stepResourceName := ""
	resourceName := args[1]
	result := make(map[string]interface{})

	// Step - flow_run_name.step_name
	if resType == "step" && strings.Contains(args[1], ".") {
		resType = "flowrun"

		splits := strings.Split(args[1], ".")
		resourceName = splits[0]
		stepResourceName = splits[1]
	}

	if err := g.Get(common.DefaultNamespace, resType, resourceName, &result); err != nil {
		return &ErrorReply{Message: fmt.Sprintf("resource (%s) (%s) not exist or get error (%s)", resType, resourceName, err)}
	}

	switch resType {
	case string(resource.FlowRunKind):
		return g.flowRun(result, stepResourceName)
	case string(resource.FlowKind):
		return g.flow(result)
	case string(resource.StepKind):
		return g.step(result)
	case string(resource.ActionKind):
		return g.action(result)
	}

	bs, err := json.Marshal(result)
	if err != nil {
		return &ErrorReply{Message: fmt.Sprintf("get resource (%s) unmarshal byte error (%s)", resType, err)}
	}
	typePointer := storage.GetResourceCoder(resType)
	obj := typePointer.(core.IObject).Clone()
	if err := core.JSONRawToResource(bs, obj); err != nil {
		return &ErrorReply{Message: fmt.Sprintf("get resource (%s) unmarshal byte error (%s)", resType, err)}
	}

	return &RawReply{Message: bs}
}

func arrayToStrings(t interface{}) []string {
	result := make([]string, 0)
	if t != nil {
		switch reflect.TypeOf(t).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(t)
			for i := 0; i < s.Len(); i++ {
				result = append(result, fmt.Sprintf("%s", s.Index(i)))
			}
		}
	}
	return result
}

func (g *Get) stepByFlowRun(result map[string]interface{}, stepResourceName string) map[string]interface{} {
	steps := get(result, "spec.steps")
	stepResult := make(map[string]interface{})
	if pa, Ok := steps.(primitive.A); Ok {
		stepsA := []interface{}(pa)
		for _, step := range stepsA {
			if stepResourceName == get(step.(map[string]interface{}), "metadata.name") {
				stepResult = step.(map[string]interface{})
			}
		}
	}
	return stepResult
}

func (g *Get) flowRun(result map[string]interface{}, stepResourceName string) Reply {
	format := NewFormat()
	if stepResourceName != "" {
		result = g.stepByFlowRun(result, stepResourceName)
		format.Header("name", "flow_run_id", "response_state", "global_variables", "data")
		format.Row(
			fmt.Sprintf("%s", get(result, "metadata.name")),
			fmt.Sprintf("%s", get(result, "spec.flow_id")),
			fmt.Sprintf("%s", get(result, "spec.response.state")),
			fmt.Sprintf("%s", get(result, "spec.global_variables")),
			fmt.Sprintf("%s", get(result, "spec.data")),
		)
	} else {
		format.Header("name", "uuid", "history_states", "global_variable")
		format.Row(
			fmt.Sprintf("%s", get(result, "metadata.name")),
			fmt.Sprintf("%s", get(result, "metadata.uuid")),
			strings.Join(arrayToStrings(get(result, "spec.history_states")), "\n"),
			fmt.Sprintf("%s", get(result, "spec.global_variable")),
		)
	}
	return &RawReply{format.Out()}
}

func (g *Get) step(result map[string]interface{}) Reply {
	format := NewFormat()
	format.Header("name", "flow_run_id", "response_state", "global_variables", "data")
	format.Row(
		fmt.Sprintf("%s", get(result, "metadata.name")),
		fmt.Sprintf("%s", get(result, "spec.flow_id")),
		fmt.Sprintf("%s", get(result, "spec.response.state")),
		fmt.Sprintf("%s", get(result, "spec.global_variables")),
		fmt.Sprintf("%s", get(result, "spec.data")),
	)
	return &RawReply{format.Out()}
}

func (g *Get) action(result map[string]interface{}) Reply {
	format := NewFormat()
	format.Header("name", "type", "version", "data")
	format.Row()
	return &RawReply{format.Out()}
}

func (g *Get) flow(result map[string]interface{}) Reply {
	format := NewFormat()
	format.Header("name", "type", "version", "data")
	format.Row()
	return &RawReply{format.Out()}
}

func (g *Get) Help() string {
	return `GET resource_type name`
}
