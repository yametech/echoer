package utils

import "encoding/json"

func UnstructuredObjectToInstanceObj(src interface{}, dst interface{}) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}
