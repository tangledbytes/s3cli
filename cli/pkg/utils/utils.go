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
func ValueSliceToInterfaceSlice(vs []reflect.Value, posthook func(reflect.Value) any) ([]interface{}, error) {
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

		if posthook != nil {
			is = append(is, posthook(v))
		}

		is = append(is, v.Interface())
	}

	return is, nil
}

// ContainsAny returns true if any of the element in sub slice is present in the main slice
func ContainsAny[T any](main []T, sub []T, comp func(v1, v2 T) bool) bool {
	for _, v1 := range main {
		for _, v2 := range sub {
			if comp(v1, v2) {
				return true
			}
		}
	}

	return false
}

func FillStruct(m map[string]interface{}, s any) error {
	for k, v := range m {
		if err := setField(s, k, v); err != nil {
			return err
		}
	}

	return nil
}

func setField(obj any, name string, value any) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("no such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("cannot set %s field value", name)
	}

	// structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	// if structFieldType != val.Type() {
	// 	return fmt.Errorf("provided value type %s didn't match obj field type %s", val.Type(), structFieldType)
	// }

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered from panic", r)
		}
	}()

	val = val.Convert(structFieldValue.Type())
	structFieldValue.Set(val)
	return nil
}
