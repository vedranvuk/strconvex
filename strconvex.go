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
	"errors"
	"fmt"
)

var (
	// ErrStrconvex is the base error of strconvex package.
	ErrStrconvex = errors.New("strconvex")

	// ErrInvalidPath is returned when specified accessor path is invalid.
	ErrInvalidPath = fmt.Errorf("%w: invalid path", ErrStrconvex)
	// ErrInvalidValue is returned when an invalid value is specified.
	ErrInvalidValue = fmt.Errorf("%w: invalid value", ErrStrconvex)
	// ErrInvalidArgument is returned when an invalid argument was specified.
	ErrInvalidArgument = fmt.Errorf("%w: invalid argument", ErrStrconvex)
	// ErrUnsupportedValue is returned when specified value is unsupported.
	ErrUnsupportedValue = fmt.Errorf("%w: unsupported value", ErrStrconvex)
	// ErrUnaddressableValue is returned when specified value is unaddressable.
	ErrUnaddressableValue = fmt.Errorf("%w: unadressable value", ErrStrconvex)
)
