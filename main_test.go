package main

import (
	"reflect"
	"testing"
)

func TestConvertPremitiveTypes(t *testing.T) {
	values := []struct {
		val any
		exp string
	}{
		// string
		{"string", "string"},
		// integer
		{int(1), "1"},
		{int8(2), "2"},
		{int16(3), "3"},
		{int32(4), "4"},
		{int64(5), "5"},
		{uint(6), "6"},
		{uint8(7), "7"},
		{uint16(8), "8"},
		{uint32(9), "9"},
		{uint64(10), "10"},
		// float
		{float32(1.23456789), "1.2345679"},
		{float64(0.1234567890123456789), "0.12345678901234568"},
		// boolean
		{true, "1"},
		{false, "0"},
	}
	for _, v := range values {
		typ := reflect.ValueOf(v.val).Type().String()
		if res := convertToStrMap(v.val); v.exp != res {
			t.Errorf("Error: converting type(%s), expected %#v, got %#v", typ, v.exp, res)
		}
	}
}

func TestConvertSliceAndArray(t *testing.T) {
	// a slice of any contains an array of int
	sa := []any{100, 200, [2]int{300, 400}}

	exp := []any{"100", "200", []any{"300", "400"}}
	res := convertToStrMap(sa)
	results, ok := res.([]any)
	if !ok {
		t.Errorf("Error: expected %#v, got %#v", exp, res)
	}
	if !reflect.DeepEqual(results, exp) {
		t.Errorf("Error: expected %#v, got %#v", exp, res)
	}
}

func TestConvertMap(t *testing.T) {
	// a slice of any contains an array of int
	m := map[string]any{
		"key1": 100,
		"key2": 200,
		"key3": map[string]any{
			"key4": true,
			"key5": false,
		},
	}

	exp := map[string]any{
		"key1": "100",
		"key2": "200",
		"key3": map[string]any{
			"key4": "1",
			"key5": "0",
		},
	}
	res := convertToStrMap(m)
	resMap, ok := res.(map[string]any)
	if !ok {
		t.Errorf("Error: expected %#v, got %#v", exp, res)
	}

	if !reflect.DeepEqual(resMap, exp) {
		t.Errorf("Error: expected %#v, got %#v", exp, res)
	}
}

func TestConvertNil(t *testing.T) {
	res := convertToStrMap(nil)
	s, ok := res.(string)
	if !ok {
		t.Errorf("Error: failed to convert nil")
	}
	if s != "" {
		t.Errorf("Error: expected %#v, got %#v", `""`, `"`+s+`"`)
	}
}

func TestConvertStruct(t *testing.T) {
	type User struct {
		No      int    `json:"no"`
		Name    string `json:"name"`
		IsAdmin bool   `json:"is_admin"`
	}
	type Group struct {
		Name  string
		Users []User
	}
	g := Group{
		Name: "Group1",
		Users: []User{
			{No: 1, Name: "User1", IsAdmin: true},
			{No: 2, Name: "User2", IsAdmin: false},
		},
	}
	exp := map[string]any{
		"Name": "Group1",
		"Users": []any{
			map[string]any{
				"no":       "1",
				"name":     "User1",
				"is_admin": "1",
			},
			map[string]any{
				"no":       "2",
				"name":     "User2",
				"is_admin": "0",
			},
		},
	}
	res := convertToStrMap(g)
	resMap, ok := res.(map[string]any)
	if !ok {
		t.Errorf("Error: expected %#v, got %#v", exp, res)
	}
	if !reflect.DeepEqual(resMap, exp) {
		t.Errorf("Error: expected %#v, got %#v", exp, res)
	}
}
