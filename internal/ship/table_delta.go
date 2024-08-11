package ship

import (
	"fmt"
	"reflect"
)

func parseTableDeltaDataInner(v reflect.Value) reflect.Value {
	if IsVariant(v) {
		v = v.Index(1)
	}

	switch v.Kind() {
	case reflect.Interface:
		return parseTableDeltaDataInner(v.Elem())
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			v.Index(i).Set(parseTableDeltaDataInner(v.Index(i)))
		}
	case reflect.Map:
		it := v.MapRange()
		for it.Next() {
			v.SetMapIndex(it.Key(), parseTableDeltaDataInner(it.Value()))
		}
	}

	return v
}

func ParseTableDeltaData(v any) (map[string]interface{}, error) {
	iface := parseTableDeltaDataInner(reflect.ValueOf(v)).Interface()
	if out, ok := iface.(map[string]interface{}); ok {
		return out, nil
	}
	return nil, fmt.Errorf("data is not an map")
}
