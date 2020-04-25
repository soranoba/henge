package henge

import (
	"fmt"
	"reflect"
)

func Copy(in interface{}, out interface{}) {
	deepCopy(reflect.ValueOf(in), reflect.ValueOf(out))
}

func deepCopy(in reflect.Value, out reflect.Value) {
	if out.Kind() == reflect.Ptr {
		out = out.Elem()
	}

	if !out.CanSet() {
		panic("henge.Copy only allows values that can be writable out")
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

		for {
			pair := ite.More()
			if pair == nil {
				break
			}
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

		for {
			pair := ite.More()
			if pair == nil {
				break
			}
			v := out.FieldByName(String(pair.Key.Interface()))
			fmt.Println(pair, v, v.Kind(), v.IsValid(), v.CanSet())
			if v.IsValid() {
				if dst := reflect.Indirect(v); dst.CanSet() {
					deepCopy(pair.Value, dst)
				} else {
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

		for idx := 0; ; idx++ {
			pair := ite.More()
			if pair == nil {
				break
			}
			if idx < out.Len() {
				deepCopy(pair.Value, out.Index(idx))
			} else {
				c := reflect.New(out.Type().Elem()).Elem()
				deepCopy(pair.Value, c)
				out.Set(reflect.Append(out, c))
			}
		}
	}
}
