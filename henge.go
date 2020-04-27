package henge

import (
	"fmt"
	"reflect"
	"strconv"
)

// ref: MapWithDepth()
const UnlimitedDepth uint = ^uint(0)

// Converts and returns interface to int.
// If it cannot convert, it returns zero.
func Int(in interface{}) (out int64) {
	inV := reflect.Indirect(reflect.ValueOf(in))
	if !inV.IsValid() {
		return
	}

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

// Converts and returns interface to uint.
// If it cannot convert, it returns zero.
func Uint(in interface{}) (out uint64) {
	inV := reflect.Indirect(reflect.ValueOf(in))
	if !inV.IsValid() {
		return
	}

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

// Converts and returns interface to float.
// If it cannot convert, it returns zero.
func Float(in interface{}) (out float64) {
	inV := reflect.Indirect(reflect.ValueOf(in))
	if !inV.IsValid() {
		return
	}

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

// Converts and returns interface to string.
// If it cannot convert, it returns empty string.
func String(in interface{}) (out string) {
	inV := reflect.Indirect(reflect.ValueOf(in))
	if !inV.IsValid() {
		return
	}

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

// Converts and returns interface to a pointer of int.
//
// It only returns nil, if nil inputted.
// Otherwise, it returns a pointer of result conversion to int.
func IntPtr(in interface{}) *int64 {
	if v := reflect.Indirect(reflect.ValueOf(in)); !v.IsValid() {
		return nil
	}

	i := Int(in)
	return &i
}

// Converts and returns interface to a pointer of uint.
//
// It only returns nil, if nil inputted.
// Otherwise, it returns a pointer of result conversion to uint.
func UintPtr(in interface{}) *uint64 {
	if v := reflect.Indirect(reflect.ValueOf(in)); !v.IsValid() {
		return nil
	}

	ui := Uint(in)
	return &ui
}

// Converts and returns interface to a pointer of float.
//
// It only returns nil, if nil inputted.
// Otherwise, it returns a pointer of result conversion to float.
func FloatPtr(in interface{}) *float64 {
	if v := reflect.Indirect(reflect.ValueOf(in)); !v.IsValid() {
		return nil
	}

	f := Float(in)
	return &f
}

// Converts and returns interface to a pointer of string.
//
// It only returns nil, if nil inputted.
// Otherwise, it returns a pointer of result conversion to string.
func StringPtr(in interface{}) *string {
	if v := reflect.Indirect(reflect.ValueOf(in)); !v.IsValid() {
		return nil
	}

	s := String(in)
	return &s
}

// Convert a struct to a map.
// It will not be expanded nested structs.
func Map(in interface{}) (out map[string]interface{}) {
	return MapWithDepth(in, 0)
}

// MapWithDepth() is the same as Map().
// It can specify how much to expand the nested struct.
// When depth is zero, it returns empty map.
func MapWithDepth(in interface{}, depth uint) (out map[string]interface{}) {
	inV := reflect.Indirect(reflect.ValueOf(in))
	if !inV.IsValid() {
		return make(map[string]interface{}, 0)
	}

	inT := inV.Type()
	if inT.Kind() != reflect.Struct {
		panic("henge.Map can only support conversion from structs")
	}

	out = make(map[string]interface{}, inT.NumField())
	for i := 0; i < inT.NumField(); i++ {
		field := inT.Field(i)
		v := reflect.Indirect(inV.FieldByName(field.Name))

		if v.Kind() == reflect.Struct && depth > 0 {
			out[field.Name] = MapWithDepth(v.Interface(), depth-1)
		} else if v.IsValid() && v.CanInterface() {
			out[field.Name] = v.Interface()
		}
	}
	return
}
