// Copyright 2020 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package strconvex is an extension to strconv package that implements a common
// interface for converting text to simple Go values using reflect to guess
// type information and strconv for actual conversion.
//
// As in strconv, only simple Go types are supported with a few minor helpful
// additions that help with compound types but have the limitation that only the
// first level is parsed and their elements or fields must be simple types.
//
// As input, standard GoValue format from the fmt package is understood.
package strconvex

import (
	"encoding"
	"errors"
	"reflect"
	"strconv"
	"strings"
)

// StringToInterface converts string in to out which must be a pointer to a Go
// value conversion compatible to data contained in string. See StringToValue
// for details.
func StringToInterface(in string, out interface{}) error {
	if out == nil {
		return ErrInvalidArgument
	}
	return stringToInterface(in, out)
}

func stringToInterface(in string, out interface{}) error {
	return stringToValue(in, reflect.Indirect(reflect.ValueOf(out)))
}

// StringToValue converts string in to out reflect.Value whose type must be
// conversion compatible to in string.
//
// Compound types are supported on first level only and must have simple types
// as elements. Simple parsing rules are defined as follows and reserve the
// comma ',', equals '=', left brace'{' and right brace '}' characters for
// interpreting the compound types. Compound values are specified as follows:
//
// Array and Slice: Values delimited by comma.
// Example: 0,1,2,3,4
//
// Map: Key=Value pairs delimited by comma.
// Example:key1=value1,key2=value2,keyN=valueN
//
// Struct: Map of values enclosed in braces.
// Example:{field1=foo,field2=42,fieldN=valueN}
//
// Invalid syntax for compound values, or Chans and Func values as input values
// will result in an error.
func StringToValue(in string, out reflect.Value) error {
	if !out.IsValid() {
		return ErrInvalidValue
	}
	return stringToValue(in, out)
}

func stringToValue(in string, out reflect.Value) error {
	bum, ok := out.Interface().(encoding.TextUnmarshaler)
	if ok {
		if err := bum.UnmarshalText([]byte(in)); err != nil {
			return err
		}
		return nil
	}
	switch out.Kind() {
	case reflect.Bool:
		return stringToBoolValue(in, out)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return stringToIntValue(in, out)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return stringToUintValue(in, out)
	case reflect.Float32:
		return stringToFloat32Value(in, out)
	case reflect.Float64:
		return stringToFloat64Value(in, out)
	case reflect.Complex64:
		return stringToComplex64Value(in, out)
	case reflect.Complex128:
		return stringToComplex128Value(in, out)
	case reflect.String:
		return stringToStringValue(in, out)
	case reflect.Array:
		return stringToArrayValue(in, out)
	case reflect.Slice:
		return stringToSliceValue(in, out)
	case reflect.Map:
		return stringToMapValue(in, out)
	case reflect.Struct:
		return stringToStructValue(in, out)
	case reflect.Ptr:
		return stringToPointerValue(in, out)
	}
	return ErrUnsupportedValue
}

// StringToBoolValue converts a string to a bool.
func StringToBoolValue(in string, out reflect.Value) error {
	if !out.IsValid() {
		return ErrInvalidValue
	}
	return stringToBoolValue(in, out)
}

// StringToBoolValue is the implementation of stringToBoolValue.
func stringToBoolValue(in string, out reflect.Value) error {
	b, err := strconv.ParseBool(in)
	if err != nil {
		return err
	}
	out.Set(reflect.ValueOf(b))
	return nil
}

// StringToIntValue converts a string to a int of any width.
func StringToIntValue(in string, out reflect.Value) error {
	if !out.IsValid() {
		return ErrInvalidValue
	}
	return stringToIntValue(in, out)
}

// StringToIntValue is the implementation of stringToIntValue.
func stringToIntValue(in string, out reflect.Value) error {
	n, err := strconv.ParseInt(in, 10, 64)
	if err != nil {
		return err
	}
	out.Set(reflect.ValueOf(n).Convert(out.Type()))
	return nil
}

// StringToUintValue converts a string to an uint of any width.
func StringToUintValue(in string, out reflect.Value) error {
	if !out.IsValid() {
		return ErrInvalidValue
	}
	return stringToUintValue(in, out)
}

func stringToUintValue(in string, out reflect.Value) error {
	n, err := strconv.ParseUint(in, 10, 64)
	if err != nil {
		return err
	}
	out.Set(reflect.ValueOf(n).Convert(out.Type()))
	return nil
}

// StringToFloat32Value converts a string to a float32.
func StringToFloat32Value(in string, out reflect.Value) error {
	if !out.IsValid() {
		return ErrInvalidValue
	}
	return stringToFloat32Value(in, out)
}

func stringToFloat32Value(in string, out reflect.Value) error {
	n, err := strconv.ParseFloat(in, 32)
	if err != nil {
		return err
	}
	out.Set(reflect.ValueOf(n).Convert(out.Type()))
	return nil
}

// StringToFloat64Value converts a string to a float64.
func StringToFloat64Value(in string, out reflect.Value) error {
	if !out.IsValid() {
		return ErrInvalidValue
	}
	return stringToFloat64Value(in, out)
}

func stringToFloat64Value(in string, out reflect.Value) error {
	n, err := strconv.ParseFloat(in, 64)
	if err != nil {
		return err
	}
	out.Set(reflect.ValueOf(n).Convert(out.Type()))
	return nil
}

// StringToComplex64Value converts a string to a complex64.
func StringToComplex64Value(in string, out reflect.Value) error {
	if !out.IsValid() {
		return ErrInvalidValue
	}
	return stringToComplex64Value(in, out)
}

func stringToComplex64Value(in string, out reflect.Value) error {
	n, err := strconv.ParseComplex(in, 64)
	if err != nil {
		return err
	}
	out.Set(reflect.ValueOf(n).Convert(out.Type()))
	return nil
}

// StringToComplex128Value converts a string to a complex128.
func StringToComplex128Value(in string, out reflect.Value) error {
	if !out.IsValid() {
		return ErrInvalidValue
	}
	return stringToComplex128Value(in, out)
}

func stringToComplex128Value(in string, out reflect.Value) error {
	n, err := strconv.ParseComplex(in, 128)
	if err != nil {
		return err
	}
	out.Set(reflect.ValueOf(n).Convert(out.Type()))
	return nil
}

// StringToStringValue converts a string to a string.
func StringToStringValue(in string, out reflect.Value) error {
	if !out.IsValid() {
		return ErrInvalidValue
	}
	return stringToStringValue(in, out)
}

func stringToStringValue(in string, out reflect.Value) error {
	out.Set(reflect.ValueOf(in))
	return nil
}

// StringToArrayValue converts a string to an array.
// String is of the form "elem1,elem2,elemN".
// Elements must be simple values.
func StringToArrayValue(in string, out reflect.Value) error {
	if !out.IsValid() {
		return ErrInvalidValue
	}
	return stringToArrayValue(in, out)
}

func stringToArrayValue(in string, out reflect.Value) error {
	v := reflect.Indirect(reflect.New(out.Type()))
	a := strings.Split(in, ",")
	for i, l := 0, out.Len(); i < l && i < len(a); i++ {
		if err := StringToValue(strings.TrimSpace(a[i]), v.Index(i)); err != nil {
			return err
		}
	}
	out.Set(v)
	return nil
}

// StringToSliceValue converts a string to a slice.
// String is of the form "elem1,elem2,elemN".
// Elements must be simple values.
func StringToSliceValue(in string, out reflect.Value) error {
	if !out.IsValid() {
		return ErrInvalidValue
	}
	return stringToSliceValue(in, out)
}

func stringToSliceValue(in string, out reflect.Value) error {
	a := strings.Split(in, ",")
	parsedval := reflect.MakeSlice(reflect.SliceOf(out.Type().Elem()), len(a), len(a))
	for i := 0; i < len(a); i++ {
		if err := StringToValue(a[i], parsedval.Index(i)); err != nil {
			return err
		}
	}
	out.Set(parsedval)
	return nil
}

// StringToMapValue converts a string to a map.
// String is of the form: "key1=val1,key2=val2,keyN=valN".
// Elements must be simple values.
func StringToMapValue(in string, out reflect.Value) error {
	if !out.IsValid() {
		return ErrInvalidValue
	}
	return stringToMapValue(in, out)
}

func stringToMapValue(in string, out reflect.Value) error {
	var maptype = reflect.MapOf(out.Type().Key(), out.Type().Elem())
	var newmap = reflect.MakeMap(maptype)
	var key, val reflect.Value
	var pair []string
	for _, s := range strings.Split(in, ",") {
		if pair = strings.Split(strings.TrimSpace(s), "="); len(pair) != 2 {
			return errors.New("strconvex: syntax error")
		}
		key = reflect.Indirect(reflect.New(maptype.Key()))
		if err := StringToValue(pair[0], key); err != nil {
			return err
		}
		val = reflect.Indirect(reflect.New(maptype.Elem()))
		if err := StringToValue(pair[1], val); err != nil {
			return err
		}
		newmap.SetMapIndex(key, val)
	}
	out.Set(newmap)
	return nil
}

// StringToStructValue converts a string to a struct.
// String is of the form: "{field1=value1,field2=value2,fieldN=valueN}"
// Fields must be simple values.
func StringToStructValue(in string, out reflect.Value) error {
	if !out.IsValid() {
		return ErrInvalidValue
	}
	return stringToStructValue(in, out)
}

func stringToStructValue(in string, out reflect.Value) error {
	var pair []string
	var field reflect.Value
	var val reflect.Value
	for _, s := range strings.Split(strings.TrimPrefix(strings.TrimSuffix(in, "}"), "{"), ",") {
		if pair = strings.Split(strings.TrimSpace(s), "="); len(pair) != 2 {
			return errors.New("strconvex: syntax error")
		}
		if field = out.FieldByName(pair[0]); !field.IsValid() {
			return errors.New("strconvex: field not found")
		}
		val = reflect.Indirect(reflect.New(field.Type()))
		if err := StringToValue(pair[1], val); err != nil {
			return err
		}
		field.Set(val)
	}
	return nil
}

// StringToPointerValue converts a string to a pointer.
func StringToPointerValue(in string, out reflect.Value) error {
	if !out.IsValid() {
		return ErrInvalidValue
	}
	return stringToPointerValue(in, out)
}

func stringToPointerValue(in string, out reflect.Value) error {
	nv := reflect.New(out.Type().Elem())
	if err := StringToValue(in, reflect.Indirect(nv)); err != nil {
		return err
	}
	out.Set(nv)
	return nil
}
