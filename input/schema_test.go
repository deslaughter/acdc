package input_test

import (
	"fmt"
	"reflect"
)

func CompareStructs(a, b any) error {
	aVal, bVal := reflect.ValueOf(a).Elem(), reflect.ValueOf(b).Elem()
	for i := 0; i < aVal.NumField(); i++ {
		if !aVal.Type().Field(i).IsExported() {
			continue
		}
		act, exp := aVal.Field(i).Interface(), bVal.Field(i).Interface()
		if !reflect.DeepEqual(act, exp) {
			return fmt.Errorf("field %s = '%v', expected '%v'", aVal.Type().Field(i).Name, act, exp)
		}
	}
	return nil
}
