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
		if v.Type().Implements(errorInterface) {
			if !v.IsNil() {
				is = append(is, map[string]interface{}{
					"error": v.Interface().(error).Error(),
				})
			} else {
				is = append(is, map[string]interface{}{
					"error": nil,
				})
			}

			continue
		}

		if !v.CanInterface() {
			return nil, fmt.Errorf("cannot convert value at idx: %d to interface", i)
		}

		is = append(is, v.Interface())

		if posthook != nil {
			is = append(is, posthook(v))
		}
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
	val := reflect.ValueOf(value)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("no such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("cannot set %s field value", name)
	}

	conv, ok := convertValue(val, structFieldValue)
	if !ok {
		return fmt.Errorf("cannot convert %s field value: type of value %q while field %q", name, val.Type(), structFieldValue.Type())
	}

	structFieldValue.Set(conv)
	return nil
}

func convertValue(from, to reflect.Value) (reflect.Value, bool) {
	if from.Type() == to.Type() {
		return from, true
	}

	if from.Type().ConvertibleTo(to.Type()) {
		return from.Convert(to.Type()), true
	}

	if to.Type() == reflect.PointerTo(from.Type()) {
		data := reflect.New(from.Type())
		data.Elem().Set(from)

		return data, true
	}

	return reflect.Value{}, false
}
