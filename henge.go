package henge

import (
	"fmt"
	"reflect"
	"strconv"
)

func Int(in interface{}) (out int64) {
	if in == nil {
		return
	}

	inV := reflect.Indirect(reflect.ValueOf(in))
	inT := inV.Type()
	outT := reflect.TypeOf(out)
	if inT.ConvertibleTo(outT) {
		out = inV.Convert(outT).Interface().(int64)
	} else if inT.Kind() == reflect.String {
		out, _ = strconv.ParseInt(inV.Interface().(string), 10, 64)
	} else if inT.Kind() == reflect.Bool {
		if inV.Interface().(bool) == true {
			out = 1
		}
	}
	return
}

func Uint(in interface{}) (out uint64) {
	if in == nil {
		return
	}

	inV := reflect.Indirect(reflect.ValueOf(in))
	inT := inV.Type()
	outT := reflect.TypeOf(out)
	if inT.ConvertibleTo(outT) {
		out = inV.Convert(outT).Interface().(uint64)
	} else if inT.Kind() == reflect.String {
		out, _ = strconv.ParseUint(inV.Interface().(string), 10, 64)
	} else if inT.Kind() == reflect.Bool {
		if inV.Interface().(bool) == true {
			out = 1
		}
	}
	return
}

func Float(in interface{}) (out float64) {
	if in == nil {
		return
	}

	inV := reflect.Indirect(reflect.ValueOf(in))
	inT := inV.Type()
	outT := reflect.TypeOf(out)
	if inT.ConvertibleTo(outT) {
		out = inV.Convert(outT).Interface().(float64)
	} else if inT.Kind() == reflect.String {
		out, _ = strconv.ParseFloat(inV.Interface().(string), 64)
	} else if inT.Kind() == reflect.Bool {
		if inV.Interface().(bool) == true {
			out = 1
		}
	}
	return
}

func String(in interface{}) (out string) {
	if in == nil {
		return
	}

	inV := reflect.Indirect(reflect.ValueOf(in))
	inT := inV.Type()
	outT := reflect.TypeOf(out)

	if inT.Kind() == reflect.Int ||
		inT.Kind() == reflect.Int8 ||
		inT.Kind() == reflect.Int16 ||
		inT.Kind() == reflect.Int32 ||
		inT.Kind() == reflect.Int64 {
		var i int64
		i = inV.Convert(reflect.TypeOf(i)).Interface().(int64)
		out = strconv.FormatInt(i, 10)
	} else if inT.Kind() == reflect.Uint ||
		inT.Kind() == reflect.Uint8 ||
		inT.Kind() == reflect.Uint16 ||
		inT.Kind() == reflect.Uint32 ||
		inT.Kind() == reflect.Uint64 {
		var ui uint64
		ui = inV.Convert(reflect.TypeOf(ui)).Interface().(uint64)
		out = strconv.FormatUint(ui, 10)
	} else if inT.Kind() == reflect.Float32 ||
		inT.Kind() == reflect.Float64 {
		var f float64
		f = inV.Convert(reflect.TypeOf(f)).Interface().(float64)
		out = fmt.Sprintf("%f", f)
	} else if inT.Kind() == reflect.Bool {
		if inV.Interface().(bool) == true {
			out = "true"
		} else {
			out = "false"
		}
	} else if inT.ConvertibleTo(outT) {
		out = inV.Convert(outT).Interface().(string)
	}
	return
}

func StringPtr(in interface{}) *string {
	s := String(in)
	return &s
}

func Map(in interface{}) (out map[string]interface{}) {
	inV := reflect.Indirect(reflect.ValueOf(in))
	inT := inV.Type()

	if inT.Kind() != reflect.Struct {
		panic("henge.Map can only support conversion from structs")
	}

	out = make(map[string]interface{}, inT.NumField())
	for i := 0; i < inT.NumField(); i++ {
		field := inT.Field(i)
		v := reflect.Indirect(inV.FieldByName(field.Name))

		if v.Kind() == reflect.Struct {
			out[field.Name] = Map(v.Interface())
		} else {
			out[field.Name] = v.Interface()
		}
	}
	return
}
