package henge

import (
	"fmt"
	"reflect"
)

type Options map[string]interface{}

// When henge.Copy is executed, afterHenge is called
// if this interface is implemented on the output type.
//
// No callback is provided before processing,
// because it copy including the private field when input and output are same type.
type AfterCallback interface {
	// src is a non-pointer type.
	AfterHenge(src interface{}, opt Options)
}

// Convert an interface to another.
// Example:
//
//   src := time.Now()
//   dst = new(time.Time)
//   Copy(src, dst)
//
// Conversion rules:
//   - Can convert between non-pointer type and pointer type.
//   - Can convert primitive type like by PHP.
//   - If the type is the same, copy it including private fields.
//   - If the type is another, copy it excluding private fields.
func Copy(in interface{}, out interface{}, opts ...Options) {
	deepCopy(reflect.ValueOf(in), reflect.ValueOf(out), opts...)
}

func deepCopy(in reflect.Value, out reflect.Value, opts ...Options) {
	if len(opts) > 1 {
		panic("opts disallow two or more thena")
	}

	if out.Kind() == reflect.Ptr {
		out = out.Elem()
	}

	if !out.CanSet() {
		panic(fmt.Sprintf("henge.Copy only allows values that can be writable out: %#v", out))
	}

	in = reflect.Indirect(in)
	if !in.IsValid() {
		return
	}

	if afterCallback, ok := out.Addr().Interface().(AfterCallback); ok {
		var opt Options
		if len(opts) == 0 {
			opt = make(Options)
		} else {
			opt = opts[0]
		}
		defer afterCallback.AfterHenge(in.Interface(), opt)
	}

	// Types that are simply converted (it also copies private fields)
	if in.Type().ConvertibleTo(out.Type()) {
		out.Set(in.Convert(out.Type()))
		return
	}

	switch out.Kind() {
	case reflect.Int:
		out.Set(reflect.ValueOf((int)(Int(in.Interface()))))
	case reflect.Int8:
		out.Set(reflect.ValueOf((int8)(Int(in.Interface()))))
	case reflect.Int16:
		out.Set(reflect.ValueOf((int16)(Int(in.Interface()))))
	case reflect.Int32:
		out.Set(reflect.ValueOf((int32)(Int(in.Interface()))))
	case reflect.Int64:
		out.Set(reflect.ValueOf((int64)(Int(in.Interface()))))
	case reflect.Uint:
		out.Set(reflect.ValueOf((uint)(Uint(in.Interface()))))
	case reflect.Uint8:
		out.Set(reflect.ValueOf((uint8)(Uint(in.Interface()))))
	case reflect.Uint16:
		out.Set(reflect.ValueOf((uint16)(Uint(in.Interface()))))
	case reflect.Uint32:
		out.Set(reflect.ValueOf((uint32)(Uint(in.Interface()))))
	case reflect.Uint64:
		out.Set(reflect.ValueOf((uint64)(Uint(in.Interface()))))
	case reflect.Float32:
		out.Set(reflect.ValueOf((float32)(Float(in.Interface()))))
	case reflect.Float64:
		out.Set(reflect.ValueOf((float64)(Float(in.Interface()))))
	case reflect.String:
		out.Set(reflect.ValueOf(String(in.Interface())))
	case reflect.Map:
		ite, ok := NewIterator(in)
		if !ok {
			return
		}

		kt := out.Type().Key()
		vt := out.Type().Elem()

		// Initialize if out is nil
		if out.IsNil() {
			out.Set(reflect.MakeMap(out.Type()))
		}

		for _, pair := range getPairs(ite) {
			kc := reflect.New(kt)
			vc := reflect.New(vt)
			deepCopy(pair.Key, kc)
			deepCopy(pair.Value, vc)
			out.SetMapIndex(pair.Key, pair.Value)
		}
	case reflect.Struct:
		ite, ok := NewIterator(in)
		if !ok {
			return
		}

		for _, pair := range getPairs(ite) {
			if pair.Value.Kind() == reflect.Ptr && pair.Value.IsNil() {
				continue
			}

			fieldName := String(pair.Key.Interface())
			structField, ok := out.Type().FieldByName(fieldName)
			if !ok || len(structField.Index) == 0 {
				continue
			}

			// NOTE: Uninitialized embeded field.
			v := out.Field(structField.Index[0])
			for v.Kind() == reflect.Ptr && v.Type().Elem().Kind() == reflect.Struct {
				if v.IsNil() {
					v.Set(reflect.New(v.Type().Elem()))
				}
				v = v.Elem()
			}

			v = out.FieldByName(fieldName)
			if v.IsValid() {
				if dst := reflect.Indirect(v); dst.CanSet() {
					deepCopy(pair.Value, dst)
				} else if v.Type().Kind() == reflect.Ptr && v.CanSet() {
					dst := reflect.New(v.Type().Elem())
					deepCopy(pair.Value, dst)
					v.Set(dst)
				}
			}
		}
	case reflect.Slice:
		fallthrough
	case reflect.Array:
		ite, ok := NewIterator(in)
		if !ok {
			return
		}

		for idx, pair := range getPairs(ite) {
			if idx < out.Len() {
				v := out.Index(idx)
				if v.Kind() == reflect.Ptr && v.IsNil() {
					v.Set(reflect.New(v.Type().Elem()))
				}
				deepCopy(pair.Value, v)
			} else if out.Kind() == reflect.Array {
				break
			} else {
				c := reflect.New(out.Type().Elem()).Elem()
				if c.Kind() == reflect.Ptr && c.IsNil() {
					c.Set(reflect.New(c.Type().Elem()))
				}
				deepCopy(pair.Value, c)
				out.Set(reflect.Append(out, c))
			}
		}
	}
}

func getPairs(ite Iterator) []*Pair {
	pairs := make([]*Pair, 0, ite.Count())
	for i := 0; ; i++ {
		pair := ite.More()
		if pair == nil {
			break
		}

		// ignore private fields
		if pair.Key.CanInterface() && pair.Value.CanInterface() {
			pairs = append(pairs, pair)
		}
	}
	return pairs
}
