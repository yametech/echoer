package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

func map2filter(m map[string]interface{}) bson.D {
	result := make(bson.D, 0)
	for key, value := range m {
		result = append(result, bson.E{Key: key, Value: value})
	}
	return result
}

func expr2labels(expr string) bson.D {
	result := bson.D{}
	switch {
	case strings.Contains(expr, ",") && strings.Contains(expr, ":"): // A:1,C:4
		for _, item := range strings.Split(expr, ",") {
			keyValue := strings.Split(item, ":")
			if len(keyValue) != 2 {
				continue
			}
			result = append(result, bson.E{Key: keyValue[0], Value: keyValue[1]})
		}

	case strings.Contains(expr, ":"): // C:4
		keyValue := strings.Split(expr, ":")
		if len(keyValue) != 2 {
			break
		}
		result = append(result, bson.E{Key: keyValue[0], Value: keyValue[1]})
	case strings.Contains(expr, ",") && strings.Contains(expr, "="): // A=1,B=4,C=1
		for _, item := range strings.Split(expr, ",") {
			keyValue := strings.Split(item, "=")
			if len(keyValue) != 2 {
				continue
			}
			result = append(result, bson.E{Key: keyValue[0], Value: keyValue[1]})
		}
	case strings.Contains(expr, "="): //C=1
		keyValue := strings.Split(expr, "=")
		if len(keyValue) != 2 {
			break
		}
		result = append(result, bson.E{Key: keyValue[0], Value: keyValue[1]})
	}

	return result
}
