package fishfish

import (
	"encoding/json"
)

// Converts a map of JSON values to a struct
// Used for WebSockets
func JSONStructToMap[T any](m map[string]interface{}) (*T, error) {
	jsonString, err := json.Marshal(m)

	if err != nil {
		return nil, err
	}

	var finalStruct T
	if err := json.Unmarshal(jsonString, &finalStruct); err != nil {
		return nil, err
	}

	return &finalStruct, nil
}
