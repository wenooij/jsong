package jsong

import (
	"fmt"
	"reflect"
	"strings"
)

var (
	booleanType = reflect.TypeOf(boolean(false))
	numType     = reflect.TypeOf(num(0))
	strType     = reflect.TypeOf(str(""))
	stringType  = reflect.TypeOf(string(""))
)

// ValueOf creates the jsong value of the input v.
//
// It performs an operation similar to, but more
// efficient than:
//
//	data, _ := json.Marshal(x)
//	v := new(any)
//	json.Unmarshal(data, v)
//
// ValueOf panics if v contains any cycles.
func ValueOf(v any) any {
	if _, ok := v.(valueInterface); ok {
		// No need to re-encode valueInterface.
		return v
	}
	return new(encoder).encode(reflect.ValueOf(v))
}

type encoder struct {
	// Avoid cycles.
	// See pkg.go.dev/encoding/json#encodeState for details.
	ptrLevel uint
	ptrSeen  map[any]struct{}
}

const startDetectingCyclesAfter = 100

func (e *encoder) encode(v reflect.Value) valueInterface {
	switch v.Kind() {
	case reflect.Bool:
		return v.Convert(booleanType).Interface().(boolean)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint,
		reflect.Float32, reflect.Float64:
		return v.Convert(numType).Interface().(valueInterface)
	case reflect.String:
		return v.Convert(strType).Interface().(valueInterface)
	case reflect.Array:
		return e.encodeArray(v)
	case reflect.Slice:
		return e.encodeArray(v)
	case reflect.Map:
		return e.encodeMap(v)
	case reflect.Struct:
		return e.encodeStruct(v)
	case reflect.Interface:
		return e.encodeInterface(v)
	case reflect.Pointer:
		return e.encodePointer(v)
	case reflect.Invalid:
		return null{}
	default:
		panic(fmt.Errorf("ValueOf: unsupported Kind: %q", v.Kind()))
	}
}

func (e *encoder) encodeInterface(v reflect.Value) valueInterface {
	if v.IsNil() {
		return nil
	}
	return e.encode(v.Elem())
}

func (e *encoder) encodePointer(v reflect.Value) valueInterface {
	if v.IsNil() {
		return nil
	}
	if e.ptrLevel++; e.ptrLevel > startDetectingCyclesAfter {
		ptr := v.Interface()
		if _, ok := e.ptrSeen[ptr]; ok {
			panic(fmt.Errorf("ValueOf: encountered a cycle via %s", v.Type()))
		}
		e.ptrSeen[ptr] = struct{}{}
		defer func() { e.ptrLevel--; delete(e.ptrSeen, ptr) }()
	}
	return e.encode(v.Elem())
}

func (e *encoder) encodeArray(v reflect.Value) valueInterface {
	if v.IsNil() {
		return array(nil)
	}
	if v.Kind() == reflect.Slice && v.Type().Elem().Kind() == reflect.Uint8 {
		// Special case for []byte.
		return v.Convert(strType).Interface().(valueInterface)
	}
	if e.ptrLevel++; e.ptrLevel > startDetectingCyclesAfter {
		ptr := struct {
			ptr interface{}
			len int
		}{v.UnsafePointer(), v.Len()}
		if _, ok := e.ptrSeen[ptr]; ok {
			panic(fmt.Errorf("ValueOf: encountered a cycle via %s", v.Type()))
		}
		e.ptrSeen[ptr] = struct{}{}
		defer func() { e.ptrLevel--; delete(e.ptrSeen, ptr) }()
	}
	res := make(array, v.Len())
	for i := 0; i < v.Len(); i++ {
		res[i] = e.encode(v.Index(i))
	}
	return res
}

func (e *encoder) encodeMap(v reflect.Value) valueInterface {
	if v.IsNil() {
		return object(nil)
	}
	if !v.Type().Key().ConvertibleTo(stringType) {
		panic(fmt.Errorf("ValueOf: map keys must be convertable to string"))
	}
	if e.ptrLevel++; e.ptrLevel > startDetectingCyclesAfter {
		ptr := v.UnsafePointer()
		if _, ok := e.ptrSeen[ptr]; ok {
			panic(fmt.Errorf("ValueOf: encountered a cycle via %s", v.Type()))
		}
		e.ptrSeen[ptr] = struct{}{}
		defer func() { e.ptrLevel--; delete(e.ptrSeen, ptr) }()
	}
	iter := v.MapRange()
	res := make(object, v.Len())
	for iter.Next() {
		k := iter.Key().Convert(stringType).Interface().(string)
		res[k] = e.encode(iter.Value())
	}
	return res
}

func (d *encoder) encodeStruct(rv reflect.Value) valueInterface {
	t := rv.Type()
	res := make(object, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name := f.Name
		var omitEmpty bool
		if jsonTag, ok := f.Tag.Lookup("json"); ok {
			omitEmpty = strings.HasSuffix(jsonTag, ",omitempty")
			name = strings.TrimSuffix(jsonTag, ",omitempty")
		}
		if name == "-" {
			continue // Skip no JSON.
		}
		v := d.encode(rv.Field(i))
		if name == "" && f.Anonymous {
			vm := v.(object)
			// Embed resulting values into the map directly
			// While avoiding collisions.
			var toDeleteKeys []string
			for k, v := range vm {
				if _, ok := res[k]; !ok {
					res[k] = v
					toDeleteKeys = append(toDeleteKeys, k)
				}
			}
			for _, k := range toDeleteKeys {
				delete(vm, k)
			}
			if len(vm) > 0 {
				// Any collision keys are stored in the original place.
				res[f.Name] = vm
			}
			continue
		}
		if !omitEmpty || !reflect.ValueOf(v).IsZero() {
			res[name] = v
		}
	}
	return res
}
