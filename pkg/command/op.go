package command

import "strings"

// data = {"a":{"b":{"c":123}}}
// get(data,"a.b.c") = 123
func get(data map[string]interface{}, path string) (value interface{}) {
	head, remain := shift(path)
	_, exist := data[head]
	if exist {
		if remain == "" {
			return data[head]
		}
		switch data[head].(type) {
		case map[string]interface{}:
			return get(data[head].(map[string]interface{}), remain)
		}
	}
	return nil
}

func shift(path string) (head string, remain string) {
	slice := strings.Split(path, ".")
	if len(slice) < 1 {
		return "", ""
	}
	if len(slice) < 2 {
		remain = ""
		head = slice[0]
		return
	}
	return slice[0], strings.Join(slice[1:], ".")
}
