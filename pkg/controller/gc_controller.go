package controller

import (
	"fmt"
	"github.com/yametech/echoer/pkg/common"
	"github.com/yametech/echoer/pkg/resource"
	"github.com/yametech/echoer/pkg/storage"
	"github.com/yametech/echoer/pkg/utils"
	"time"
)

type GCController struct {
	storage.IStorage
}

func GC(iStorage storage.IStorage) {
	gcController := GCController{
		iStorage,
	}
	for {
		time.Sleep(time.Hour * 1)
		gcController.gcFlowRun()
	}
}

func (g *GCController) gcFlowRun() {
	nTime := time.Now().Add(-time.Hour * 24).Unix()
	data, err := g.List(common.DefaultNamespace, common.FlowRunCollection, "")
	if err != nil {
		fmt.Printf("gc list flowrun error %s", err)
		return
	}
	if len(data) > 10000 {
		count := len(data) - 10000
		for idx, flowRun := range data {
			this := &resource.FlowRun{}
			err := utils.UnstructuredObjectToInstanceObj(flowRun, this)
			if err != nil {
				fmt.Printf("gc flowrun unmarshal error %s", err)
				return
			}
			if idx < count {
				g.cleanFlowRun(this, true)
				continue
			}
			if this.Version < nTime {
				g.cleanFlowRun(this, false)
			}
		}
	} else {
		for _, flowRun := range data {
			this := &resource.FlowRun{}
			err := utils.UnstructuredObjectToInstanceObj(flowRun, this)
			if err != nil {
				fmt.Printf("gc flowrun unmarshal error %s", err)
				return
			}
			if this.Version < nTime {
				g.cleanFlowRun(this, false)
			}
		}
	}
}

func (g *GCController) cleanFlowRun(this *resource.FlowRun, forceDelete bool) {
	err := g.cleanStep(this, forceDelete)
	if err != nil {
		fmt.Printf("gc flowrun step delete error %s", err)
		return
	}
	err = g.Delete(common.DefaultNamespace, common.FlowRunCollection, this.Name)
	if err != nil {
		fmt.Printf("gc flowrun flow delete error %s", err)
	}
}

func (g *GCController) cleanStep(this *resource.FlowRun, forceDelete bool) error {
	for _, step := range this.Spec.Steps {
		if step.Spec.Done != true && forceDelete != true {
			return fmt.Errorf("step not done")
		}
		err := g.Delete(common.DefaultNamespace, common.Step, step.Name)
		if err != nil {
			return err
		}
	}
	return nil
}
