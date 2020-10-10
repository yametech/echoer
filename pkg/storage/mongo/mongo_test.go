package mongo

import (
	"testing"

	"github.com/yametech/echoer/pkg/common"
	"github.com/yametech/echoer/pkg/core"
	"github.com/yametech/echoer/pkg/resource"
	"github.com/yametech/echoer/pkg/storage"
	"github.com/yametech/echoer/pkg/storage/gtm"
)

var _ core.IObject = &TestResource{}

type TestResourceSpec struct{}

type TestResource struct {
	// Metadata default IObject Metadata
	core.Metadata `json:"metadata"`
	// Spec default TestResourceSpec Spec
	Spec TestResourceSpec `json:"spec"`
}

func (a *TestResource) Clone() core.IObject {
	result := &TestResource{}
	core.Clone(a, result)
	return result
}

// To go_test this code you need to use mongodb
/*
	docker run -itd --name mongo --net=host mongo mongod --replSet rs0
	docker exec -ti mongo mongo
	use admin;
	var cfg = {
		"_id": "rs0",
		"protocolVersion": 1,
		"members": [
			{
				"_id": 0,
				"host": "172.16.241.131:27017"
			},
		]
	};
	rs.initiate(cfg, { force: true });
	rs.reconfig(cfg, { force: true });
*/

const testIp = "172.16.241.131:27017"

// To go_test this code you need to use mongodb
func TestMongo_Apply(t *testing.T) {
	client, err := NewMongo("mongodb://" + testIp + "/admin")
	if err != nil {
		t.Fatal("open client error")
	}
	defer client.Close()
	testResource := &TestResource{
		Metadata: core.Metadata{
			Name:    "test_name",
			Kind:    core.Kind("test_resource_kind"),
			Version: 0,
			Labels:  map[string]interface{}{"who": "iam"},
		},
		Spec: TestResourceSpec{},
	}
	if _, _, err := client.Apply("default", "test_resource_kind", "test_name", testResource); err != nil {
		t.Fatal(err)
	}

	testResource.Metadata.Name = "test_name1"
	if _, _, err := client.Apply("default", "test_resource_kind", "test_name", testResource); err != nil {
		t.Fatal(err)
	}

}

func TestMongo_Watch(t *testing.T) {
	client, err := NewMongo("mongodb://" + testIp + "/admin")
	if err != nil {
		t.Fatal("open client error")
	}
	defer client.Close()

	testResource := &TestResource{
		Metadata: core.Metadata{
			Name:    "test_name",
			Kind:    core.Kind("test_resource_kind"),
			Version: 0,
			Labels:  map[string]interface{}{"who": "iam"},
		},
		Spec: TestResourceSpec{},
	}
	if _, _, err := client.Apply("default", "test_resource_kind", "test_name", testResource); err != nil {
		t.Fatal(err)
	}

	testResource.Metadata.Name = "test_name1"
	testResource.Metadata.Version = 1
	if _, _, err := client.Apply("default", "test_resource_kind", "test_name", testResource); err != nil {
		t.Fatal(err)
	}

	chains := storage.NewWatchChan()
	if err := client.Watch("default", "test_resource_kind", 0, chains); err != nil {
		t.Fatal(err)
	}

	item, ok := <-chains.ResultChan
	if !ok {
		t.Fatal("watch item not ok")
	}

	testResourceItem := &TestResource{}
	if err := core.ObjectToResource(item, testResourceItem); err != nil {
		t.Fatal(err)
	}

	if !(testResourceItem.GetResourceVersion() > 0) {
		t.Fatal("expected version failed")
	}
	//chains.Close()
}

var _ storage.Coder = &actionCoderImpl{}

type actionCoderImpl struct{}

func (c *actionCoderImpl) Decode(op *gtm.Op) (core.IObject, error) {
	action := &resource.Action{}
	if err := core.ObjectToResource(op.Data, action); err != nil {
		return nil, err
	}
	return action, nil
}

func TestMongo_Watch2(t *testing.T) {
	client, err := NewMongo("mongodb://" + testIp + "/admin")
	if err != nil {
		t.Fatal("open client error")
	}
	defer client.Close()

	action := &resource.Action{
		Metadata: core.Metadata{
			Name:    "test_action",
			Kind:    core.Kind("action"),
			Version: 0,
			Labels:  map[string]interface{}{"who": "iam"},
		},
		Spec: resource.ActionSpec{},
	}
	if _, _, err := client.Apply(common.DefaultNamespace, common.ActionCollection, action.GetName(), action); err != nil {
		t.Fatal(err)
	}

	action.Metadata.Name = "test_action1"
	action.Metadata.Version = 1
	if _, _, err := client.Apply(common.DefaultNamespace, common.ActionCollection, action.GetName(), action); err != nil {
		t.Fatal(err)
	}

	watcher := storage.NewWatch(&actionCoderImpl{})
	client.Watch2(common.DefaultNamespace, common.ActionCollection, 0, watcher)

	item, ok := <-watcher.ResultChan()
	if !ok {
		t.Fatal("watch item not ok")
	}

	if !(item.GetResourceVersion() > 0) {
		t.Fatal("expected version failed")
	}
}
