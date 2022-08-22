package util

import "encoding/json"

func StringToMap(encoded string, dest *map[string]any) error {
	return json.Unmarshal([]byte(encoded), dest)
}
