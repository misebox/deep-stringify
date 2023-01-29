package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

func convertToStrMap(data any) any {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return ""
		}
		return convertToStrMap(val.Elem())
	} else if val.Kind() == reflect.Struct {
		m := map[string]any{}
		for i := 0; i < val.NumField(); i++ {
			vf := val.Field(i)
			tf := val.Type().Field(i)
			key := tf.Name
			tag := tf.Tag.Get("json")
			if tag != "" {
				key = tag
			}
			m[key] = convertToStrMap(vf.Interface())
		}
		return m
	} else if val.Kind() == reflect.Map {
		m := map[string]any{}
		for _, k := range val.MapKeys() {
			m[k.String()] = convertToStrMap(val.MapIndex(k).Interface())
		}
		return m
	} else if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		items := []any{}
		for i := 0; i < val.Len(); i++ {
			item := convertToStrMap(val.Index(i).Interface())
			items = append(items, item)
		}
		return items
	}

	switch v := data.(type) {
	case string:
		return v
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		if v {
			return "1"
		} else {
			return "0"
		}
	case nil:
		return ""
	case []any:
		for i, e := range v {
			v[i] = convertToStrMap(e)
		}
		return v
	case map[string]any:
		for k, e := range v {
			v[k] = convertToStrMap(e)
		}
		return v
	default:
		return v
	}
}
func main() {

	type User struct {
		Name    string
		Age     uint
		IsAdmin bool
	}
	type Response struct {
		Users  []User
		Length int
		Status string
		Nil    *int
	}
	resp := Response{
		[]User{
			{"Alice", 60, true},
			{"Bob", 60, true},
			{"Carol", 50, false},
			{"Dave", 60, false},
		},
		4,
		"OK",
		nil,
	}
	fmt.Println("==== normal JSON conversion ====")
	{
		fmt.Printf("%#v\n", resp)
		bytes, err := json.Marshal(resp)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(string(bytes))
	}

	fmt.Println("==== force string JSON conversion ====")
	{
		strMappedJSON := convertToStrMap(resp)
		fmt.Printf("%#v\n", strMappedJSON)
		bytes, err := json.Marshal(strMappedJSON)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(string(bytes))
	}
}
