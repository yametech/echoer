package core

import (
	"encoding/json"
)

func ObjectToResource(data map[string]interface{}, obj IObject) error {
	bs, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bs, obj); err != nil {
		return err
	}
	return nil
}

func ObjectToMap(obj IObject) (map[string]interface{}, error) {
	var data map[string]interface{}
	bs, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bs, &data); err != nil {
		return nil, err
	}
	return data, err
}

func JSONRawToResource(raw []byte, obj IObject) error {
	if err := json.Unmarshal(raw, obj); err != nil {
		return err
	}
	return nil
}
