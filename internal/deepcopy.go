package internal

import (
	"fmt"
	"reflect"
	"time"
)

// DeepCopy only supports for primitive types (int, int8, int16, int32, int64,
// uint, uint8, uint16, uint32, uint64, float32, float64, string, bool), and
// slices[any] and map[string]any

func DeepCopyPrimitive(v any) any {
	if v == nil {
		return nil
	}

	switch x := v.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64,
		bool, string, time.Time:
		return x
	case []any:
		var n = make([]any, len(x))
		for i, e := range x {
			n[i] = DeepCopyPrimitive(e)
		}
		return n
	case []byte:
		var n = make([]byte, len(x))
		copy(n, x)
		return n
	case []rune:
		var n = make([]rune, len(x))
		copy(n, x)
		return n
	case []int:
		var n = make([]int, len(x))
		copy(n, x)
		return n
	case []float64:
		var n = make([]float64, len(x))
		copy(n, x)
		return n
	case map[string]any:
		var n = make(map[string]any, len(x))
		for k, e := range x {
			n[k] = DeepCopyPrimitive(e)
		}
		return n
	case map[string]int:
		var n = make(map[string]int, len(x))
		for k, e := range x {
			n[k] = e
		}
		return n
	case map[string]float64:
		var n = make(map[string]float64, len(x))
		for k, e := range x {
			n[k] = e
		}
		return n
	}

	var rv = reflect.ValueOf(v)
	var rt = rv.Type()
	switch rt.Kind() {
	case reflect.Array:
		var n = reflect.New(rt).Elem()
		for i := 0; i < rv.Len(); i++ {
			n.Index(i).Set(reflect.ValueOf(DeepCopyPrimitive(rv.Index(i).Interface())))
		}
		return n.Interface()
	case reflect.Slice:
		var n = reflect.MakeSlice(rt, rv.Len(), rv.Len())
		for i := 0; i < rv.Len(); i++ {
			n.Index(i).Set(reflect.ValueOf(DeepCopyPrimitive(rv.Index(i).Interface())))
		}
		return n.Interface()
	case reflect.Map:
		var n = reflect.MakeMap(rt)
		for _, k := range rv.MapKeys() {
			nk := DeepCopyPrimitive(k.Interface())
			nv := DeepCopyPrimitive(rv.MapIndex(k).Interface())
			n.SetMapIndex(reflect.ValueOf(nk), reflect.ValueOf(nv))
		}
		return n.Interface()
	}

	panic(fmt.Errorf("DeepCopy: unsupported type %T", v))
}

func reflectCopy(v any) any {
	var rv = reflect.ValueOf(v)
	var rt = rv.Type()
	switch rt.Kind() {
	case reflect.Struct:
		var n = reflect.New(rt)
		for i := 0; i < rv.NumField(); i++ {
			var nfv = n.Field(i)
			if !nfv.CanSet() {
				continue
			}
			var vfv = rv.Field(i)
			nfv.Set(reflect.ValueOf(DeepCopyPrimitive(vfv.Interface())))
		}
	case reflect.Array:
		var n = reflect.New(rt)
		for i := 0; i < rv.Len(); i++ {
			n.Index(i).Set(reflect.ValueOf(DeepCopyPrimitive(rv.Index(i).Interface())))
		}
		return n.Interface()
	case reflect.Slice:
		var n = reflect.MakeSlice(rt, rv.Len(), rv.Len())
		for i := 0; i < rv.Len(); i++ {
			n.Index(i).Set(reflect.ValueOf(DeepCopyPrimitive(rv.Index(i).Interface())))
		}
		return n.Interface()
	case reflect.Map:
		var n = reflect.MakeMap(rt)
		for _, k := range rv.MapKeys() {
			nk := DeepCopyPrimitive(k.Interface())
			nv := DeepCopyPrimitive(rv.MapIndex(k).Interface())
			n.SetMapIndex(reflect.ValueOf(nk), reflect.ValueOf(nv))
		}
		return n.Interface()
	}

	panic(fmt.Errorf("DeepCopy: unsupported type %T", v))
}
