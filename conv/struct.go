package conv

import (
	"reflect"
	"strings"
)

func StructToMap(obj interface{}) map[string]interface{} {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}
	t := v.Type()
	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Tag.Get("json")
		parts := strings.SplitN(name, ",", 2)
		if len(parts) > 0 {
			name = parts[0]
		}
		if name == "" {
			name = t.Field(i).Name
		}
		data[name] = v.Field(i).Interface()
	}
	return data
}

func StructAssign(binding interface{}, value interface{}, args ...bool) {
	argLen := len(args)
	bVal := reflect.ValueOf(binding).Elem()
	vVal := reflect.ValueOf(value).Elem()
	vTypeOfT := vVal.Type()
	for i := 0; i < vVal.NumField(); i++ {
		name := vTypeOfT.Field(i).Name
		if ok := bVal.FieldByName(name).IsValid(); ok {
			if argLen > 0 && !args[0] && vVal.Field(i).IsZero() {
				continue
			}

			if bVal.FieldByName(name).Kind() != vVal.Field(i).Kind() {
				continue
			}

			if vVal.Field(i).Kind() == reflect.Ptr {
				bVal.FieldByName(name).Set(reflect.ValueOf(vVal.Field(i).Elem().Interface()))
			} else {
				bVal.FieldByName(name).Set(reflect.ValueOf(vVal.Field(i).Interface()))
			}
		}
	}
}

func MergeMap(src, dst map[string]any) {
	if dst == nil {
		dst = make(map[string]any)
	}
	if src == nil {
		src = make(map[string]any)
	}
	for k, v := range src {
		dst[k] = v
	}
}
