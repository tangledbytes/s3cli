package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// ParseJSONToMapStringInterface parses a JSON string to a map[string]interface{}
func ParseJSONToMapStringInterface(data string) (map[string]interface{}, error) {
	mp := map[string]interface{}{}

	if err := json.Unmarshal([]byte(data), &mp); err != nil {
		return nil, err
	}

	return mp, nil
}

// ParseJSONToMapStringString parses a JSON string to a map[string]string
func ParseJSONToMapStringString(data string) (map[string]string, error) {
	mp := map[string]string{}

	if err := json.Unmarshal([]byte(data), &mp); err != nil {
		return nil, err
	}

	return mp, nil
}

// AnyToAny tries to converts any type to any type by marshaling and unmarshaling
func AnyToAny(i1, i2 any) error {
	byt, err := json.Marshal(i1)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	err = json.Unmarshal(byt, i2)
	if err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	return nil
}

// ValueSliceToInterfaceSlice tries to converts a slice of reflect.Value to a
// slice of interface{}
func ValueSliceToInterfaceSlice(vs []reflect.Value) ([]interface{}, error) {
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()

	is := make([]interface{}, 0)
	for i, v := range vs {
		if v.Type().Implements(errorInterface) && !v.IsNil() {
			is = append(is, map[string]interface{}{
				"error": v.Interface().(error).Error(),
			})

			continue
		}

		if !v.CanInterface() {
			return nil, fmt.Errorf("cannot convert value at idx: %d to interface", i)
		}

		is = append(is, v.Interface())
	}

	return is, nil
}
