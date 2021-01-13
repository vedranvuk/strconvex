// Copyright 2020 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package strconvex

import (
	"fmt"
	"reflect"
	"testing"
)

type PathTestResult struct {
	Element string
	Token
}

type PathTest struct {
	Name     string
	TestPath string
	Results  []PathTestResult
}

func TestPath(t *testing.T) {

	var tests = []PathTest{
		{
			Name:     "EmptyPath",
			TestPath: "",
			Results:  []PathTestResult{},
		},
		{
			Name:     "NameToken",
			TestPath: "Field",
			Results: []PathTestResult{
				{
					Element: "Field",
					Token:   NameToken,
				},
			},
		},
		{
			Name:     "KeyToken",
			TestPath: "[Key]",
			Results: []PathTestResult{
				{
					Element: "[Key]",
					Token:   KeyToken,
				},
			},
		},
		{
			Name:     "KeyedNameToken",
			TestPath: "Field[Key]",
			Results: []PathTestResult{
				{
					Element: "Field[Key]",
					Token:   KeyedNameToken,
				},
			},
		},
		{
			Name:     "TwoNames",
			TestPath: "Field1.Field2",
			Results: []PathTestResult{
				{
					Element: "Field1",
					Token:   NameToken,
				},
				{
					Element: "Field2",
					Token:   NameToken,
				},
			},
		},
		{
			Name:     "DotTwoNames",
			TestPath: ".Field1.Field2",
			Results: []PathTestResult{
				{
					Element: "Field1",
					Token:   NameToken,
				},
				{
					Element: "Field2",
					Token:   NameToken,
				},
			},
		},
		{
			Name:     "TwoNamesDot",
			TestPath: "Field1.Field2.",
			Results: []PathTestResult{
				{
					Element: "Field1",
					Token:   NameToken,
				},
				{
					Element: "Field2",
					Token:   NameToken,
				},
				{
					Element: "",
					Token:   InvalidToken,
				},
			},
		},
		{
			Name:     "KeyNotClosed",
			TestPath: "[",
			Results: []PathTestResult{
				{
					Element: "",
					Token:   InvalidToken,
				},
			},
		},
		{
			Name:     "KeyOpenTwice",
			TestPath: "[[",
			Results: []PathTestResult{
				{
					Element: "",
					Token:   InvalidToken,
				},
			},
		},
		{
			Name:     "KeyNotOpened",
			TestPath: "]",
			Results: []PathTestResult{
				{
					Element: "",
					Token:   InvalidToken,
				},
			},
		},
		{
			Name:     "KeyCloseTwice",
			TestPath: "]]",
			Results: []PathTestResult{
				{
					Element: "",
					Token:   InvalidToken,
				},
			},
		},
		{
			Name:     "EmptyKey",
			TestPath: "[]",
			Results: []PathTestResult{
				{
					Element: "",
					Token:   InvalidToken,
				},
			},
		},
		{
			Name:     "NothingAfterDot",
			TestPath: ".",
			Results: []PathTestResult{
				{
					Element: "",
					Token:   InvalidToken,
				},
			},
		},
		{
			Name:     "DoubleDot",
			TestPath: "..",
			Results: []PathTestResult{
				{
					Element: "",
					Token:   InvalidToken,
				},
			},
		},
		{
			Name:     "SliceOfMapToStruct",
			TestPath: "Slice[1][Key].Field",
			Results: []PathTestResult{
				{
					Element: "Slice[1]",
					Token:   KeyedNameToken,
				},
				{
					Element: "[Key]",
					Token:   KeyToken,
				},
				{
					Element: "Field",
					Token:   NameToken,
				},
			},
		},
	}

	for i := 0; i < len(tests); i++ {
		t.Run(tests[i].Name, func(t *testing.T) {
			var path = Parse(tests[i].TestPath)
			var elem string
			var token Token
			var results []PathTestResult
			for {
				elem, token = path.Next()
				if token == NoToken {
					break
				}
				results = append(results, PathTestResult{elem, token})
				if token == InvalidToken {
					break
				}
			}
			if len(results) != len(tests[i].Results) {
				fmt.Println(results)
				t.Fatalf("Result count missmatch, want '%d', got '%d'", len(tests[i].Results), len(results))
			}
			for j := 0; j < len(tests[i].Results); j++ {
				if tests[i].Results[j].Element != results[j].Element {
					t.Fatalf("Element missmatch, want '%s', got '%s'", tests[i].Results[j].Element, results[j].Element)
				}
				if tests[i].Results[j].Token != results[j].Token {
					t.Fatalf("Token missmatch, want '%s', got '%s'", tests[i].Results[j].Token, results[j].Token)
				}
			}
		})
	}
}

func BenchmarkPath(b *testing.B) {
	const testpath = "A.B[C].D[E][F].G[H][I][J]"
	var parsers = make([]*Path, 0, b.N)
	for i := 0; i < b.N; i++ {
		parsers = append(parsers, Parse(testpath))
	}
	b.ResetTimer()
	var token Token
	for i := 0; i < b.N; i++ {
		for {
			if _, token = parsers[i].Next(); token == NoToken {
				break
			}
		}
	}
}

func TestParseElement(t *testing.T) {
	var name, key string
	var err error
	if name, key, err = ParseElement("", InvalidToken); err != ErrInvalidArgument {
		t.Fatal("ParseElement failed.")
	}
	if name, key, err = ParseElement("", NoToken); err != ErrInvalidArgument {
		t.Fatal("ParseElement failed.")
	}
	if name, key, err = ParseElement("", NameToken); err != ErrInvalidArgument {
		t.Fatal("ParseElement failed.")
	}
	if name, key, err = ParseElement("Name", NameToken); name != "Name" || key != "" || err != nil {
		t.Fatal(name, key, err)
	}
	if name, key, err = ParseElement("Name[Key]", KeyedNameToken); name != "Name" || key != "Key" || err != nil {
		t.Fatal(name, key, err)
	}
}

type Child struct {
	Bool   bool
	Int    int
	String string
}

type Root struct {
	Child
	Array [5]Child
	Slice []Child
	Map   map[string]Child
}

func getData() *Root {
	return &Root{
		Array: [5]Child{
			{
				Bool:   true,
				Int:    1,
				String: "One",
			},
			{
				Bool:   false,
				Int:    2,
				String: "Two",
			},
			{
				Bool:   true,
				Int:    3,
				String: "Three",
			},
			{
				Bool:   false,
				Int:    4,
				String: "Four",
			},
			{
				Bool:   true,
				Int:    5,
				String: "Five",
			},
		},
		Slice: []Child{
			{
				Bool:   true,
				Int:    1,
				String: "One",
			},
			{
				Bool:   false,
				Int:    2,
				String: "Two",
			},
			{
				Bool:   true,
				Int:    3,
				String: "Three",
			},
			{
				Bool:   false,
				Int:    4,
				String: "Four",
			},
			{
				Bool:   true,
				Int:    5,
				String: "Five",
			},
		},
		Map: map[string]Child{
			"One": {
				Bool:   true,
				Int:    1,
				String: "One",
			},
			"Two": {
				Bool:   false,
				Int:    2,
				String: "Two",
			},
			"Three": {
				Bool:   true,
				Int:    3,
				String: "Three",
			},
			"Four": {
				Bool:   false,
				Int:    4,
				String: "Four",
			},
			"Five": {
				Bool:   true,
				Int:    5,
				String: "Five",
			},
		},
	}
}

func TestFind(t *testing.T) {
	var val reflect.Value
	var err error
	if val, err = Find("Map[Three].String", getData()); err != nil {
		t.Fatal(err)
	}
	if val.String() != "Three" {
		t.Fatal("Find failed.")
	}
}

func TestGet(t *testing.T) {
	var intf interface{}
	var err error
	if intf, err = Get("Map[Three].String", getData()); err != nil {
		t.Fatal(err)
	}
	var s string
	var ok bool
	if s, ok = intf.(string); !ok {
		t.Fatal("Get failed.")
	}
	if s != "Three" {
		t.Fatal("Get failed.")
	}
}

func TestSet(t *testing.T) {
	var data = getData()
	var err error
	if err = Set("Slice[3].String", "Foo", data); err != nil {
		t.Fatal(err)
	}
	var intf interface{}
	if intf, err = Get("Slice[3].String", data); err != nil {
		t.Fatal(err)
	}
	var s string
	var ok bool
	if s, ok = intf.(string); !ok {
		t.Fatal("Set failed.")
	}
	if s != "Foo" {
		t.Fatal("Set failed.")
	}
}