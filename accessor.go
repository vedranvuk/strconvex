// Copyright 2020 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Accessor implements reflect based functionality for inspection and
// modification of values inside structured Go values using a path-like syntax.

package strconvex

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Token defines a path token type.
type Token int

const (
	// InvalidToken marks an invalid token.
	InvalidToken Token = iota
	// NoToken marks no token found.
	NoToken
	// NameToken marks a name token.
	// A Path element that specifies a struct field.
	NameToken
	// KeyToken marks a key-only token.
	// Slice or Array index or a Map key.
	KeyToken
	// KeyedNameToken marks a name token with a key suffix.
	// Slice, Array or a Map field name with an index or key specifier.
	KeyedNameToken
)

// String implements stringer on Token.
func (t Token) String() (s string) {
	switch t {
	case InvalidToken:
		s = "InvalidToken"
	case NoToken:
		s = "NoToken"
	case NameToken:
		s = "NameToken"
	case KeyToken:
		s = "KeyToken"
	case KeyedNameToken:
		s = "KeyedNameToken"
	}
	return
}

// Parse returns a new Path parser.
func Parse(path string) *Path {
	return &Path{
		path:   path,
		length: len(path),
	}
}

// Path is a path element iterator.
type Path struct {
	path    string
	current int
	length  int
}

// Next returns next token in Path.
// Returns an empty string and NoToken if there are no tokens left.
// Returns an empty string and InvalidToken if the token is invalid.
func (p *Path) Next() (element string, token Token) {
	var i int
	for i = p.current; i < p.length; i++ {
		switch p.path[i] {
		case '[':
			switch token {
			case InvalidToken:
				token = KeyToken
			case NameToken:
				token = KeyedNameToken
			case KeyToken, KeyedNameToken:
				return "", InvalidToken
			}
		case ']':
			switch token {
			case InvalidToken, NameToken:
				return "", InvalidToken
			}
			element = p.path[p.current : i+1]
			if !validElement(&element) {
				return "", InvalidToken
			}
			p.current = i + 1
			return
		case '.':
			switch token {
			case InvalidToken:
				token = NameToken
				p.current = i + 1
				continue
			case KeyToken, KeyedNameToken:
				return "", InvalidToken
			case NameToken:
				element = p.path[p.current:i]
				if !validElement(&element) {
					return "", InvalidToken
				}
				p.current = i
				return
			}
			element = p.path[p.current : i+1]
			p.current = i + 1
			return
		default:
			if token == InvalidToken {
				token = NameToken
			}
		}
	}
	switch token {
	case InvalidToken:
		return "", NoToken
	case KeyToken, KeyedNameToken:
		return "", InvalidToken
	case NameToken:
		element = p.path[p.current:i]
		if !validElement(&element) {
			return "", InvalidToken
		}
		p.current = i
		return
	}
	return "", NoToken
}

// validElement returns if element is a valid element.
// Element is valid if its length is not zero and if its key length is not zero,
// if present.
func validElement(element *string) bool {
	if len(*element) == 0 {
		return false
	}
	if len(*element) >= 2 {
		if (*element)[len(*element)-2:] == "[]" {
			return false
		}
	}
	return true
}

// ParseElement parses element string into:
// name and key if token is KeyedNameToken.
// name if token is NameToken.
// key if token is KeyToken.
// If element is an empty string or token is not one of above listed tokens
// ParseElement returns empty name and key and ErrInvalidArgument.
func ParseElement(element string, token Token) (name, key string, err error) {
	if element == "" {
		return "", "", ErrInvalidArgument
	}
	switch token {
	case InvalidToken, NoToken:
		err = ErrInvalidArgument
	case NameToken:
		return element, "", nil
	case KeyToken:
		if key, err = strconv.Unquote(
			strings.TrimPrefix(strings.TrimSuffix(element, "]"), "["),
		); err != nil {
			return "", "", err
		}
	case KeyedNameToken:
		var a = strings.Split(element, "[")
		if len(a) != 2 {
			return "", "", ErrInvalidPath
		}
		name = a[0]
		key = strings.TrimSuffix(a[1], "]")
	}
	return
}

// Find searches for a Go value in a compound Go value specified by root by
// specified path and returns it as a reflect Value or an error.
//
// Simple go values as root are supported too but the path must be empty.
//
//
// Path syntax is as follows:
// Maps: Name[Key]
// Slices and Arrays: Name[Index]
// Struct fields: Name
// Strings must be quoted when used as keys. Quotes must be escaped.
// Elements in hierarchy are dot separated.
//
// For example:
//
// Access second element in root struct value field named "Slice":
//  Slice[1]
//
// Access struct field named "Age" in a map[string]struct entry "Example":
//  [Example].Age
//
//
func Find(path string, root interface{}) (reflect.Value, error) {
	if path == "" {
		return reflect.Value{}, ErrInvalidPath
	}
	if root == nil {
		return reflect.Value{}, ErrInvalidArgument
	}
	var current = reflect.Indirect(reflect.ValueOf(root))
	var element, name, key string
	var parser = Parse(path)
	var token Token
	var err error
loop:
	for {
		switch element, token = parser.Next(); token {
		case InvalidToken:
			return reflect.Value{}, ErrInvalidPath
		case NoToken:
			break loop
		case NameToken:
			if current.Kind() != reflect.Struct {
				return reflect.Value{}, ErrInvalidPath
			}
			current = current.FieldByName(element)
		case KeyToken:
			if _, key, err = ParseElement(element, token); err != nil {
				return reflect.Value{}, err
			}
			return valueByKey(current, key)
		case KeyedNameToken:
			if current.Kind() != reflect.Struct {
				return reflect.Value{}, ErrInvalidPath
			}
			if name, key, err = ParseElement(element, token); err != nil {
				return reflect.Value{}, err
			}
			if current = current.FieldByName(name); !current.IsValid() {
				return reflect.Value{}, ErrInvalidPath
			}
			return valueByKey(current, key)
		}
	}
	return current, nil
}

// valueByKey retrieves an Array or Slice element or a map key by specified key
// from value which must be an Array, Slice or Map. Key must be convertible to
// an integer index if value is an Array or Slice and must be in range and must
// be convertible to key type of Map if value is a map.
// If an error occurs it is returned with an invalid/zero value.
func valueByKey(value reflect.Value, key string) (reflect.Value, error) {
	var i int
	var err error
	var mapkey reflect.Value
	switch value.Kind() {
	case reflect.Array, reflect.Slice:
		if i, err = strconv.Atoi(key); err != nil {
			return reflect.Value{}, fmt.Errorf("%w: element to index: %v", ErrInvalidPath, err)
		}
		if i >= value.Len() {
			return reflect.Value{}, errors.New("index out of range")
		}
		value = value.Index(i)
	case reflect.Map:
		mapkey = reflect.Indirect(reflect.New(value.Type().Key()))
		if err = StringToValue(key, mapkey); err != nil {
			return reflect.Value{}, fmt.Errorf("%w: key to value: %v", ErrInvalidPath, err)
		}
		value = value.MapIndex(mapkey)
	default:
		return reflect.Value{}, ErrInvalidPath
	}
	return value, nil
}

// MustFind is like Find but panics on error.
func MustFind(path string, root interface{}) (v reflect.Value) {
	var err error
	if v, err = Find(path, root); err != nil {
		panic(err)
	}
	return
}

// Get retrieves a Go value from a Go compound value by path as an interface or
// returns an error.
func Get(path string, root interface{}) (interface{}, error) {
	var val reflect.Value
	var err error
	if val, err = Find(path, root); err != nil {
		return nil, err
	}
	return val.Interface(), nil
}

// MustGet is like Get but panics on error.
func MustGet(path string, root interface{}) interface{} {
	return nil
}

func Set(path, value string, root interface{}) error {
	return nil
}

// MustSet is like Set but panics on error.
func MustSet(path, value string, root interface{}) {

}
