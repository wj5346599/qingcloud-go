package qingcloud

import (
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

func TransfomRequestToParams(a interface{}) Params {
	return convertITypeToParams(reflect.ValueOf(a).Elem())
}

func convertName(s string, number ...string) string {
	var result string
	var words []string
	var lastPos int
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if i > 0 && unicode.IsUpper(rs[i]) {
			words = append(words, s[lastPos:i])
			lastPos = i
		}
	}
	if s[lastPos:] != "" {
		words = append(words, s[lastPos:])
	}

	for k, word := range words {
		if k > 0 {
			result += "_"
		}
		result += strings.ToLower(word)
	}
	if len(number) == 1 {
		result = strings.Replace(result, "_n_", "."+number[0]+".", -1)
		// 只有在最后的时候才需要替换
		if strings.Index(result, "_n") == len(result)-2 {
			result = strings.Replace(result, "_n", "."+number[0], -1)
		}

	}
	return result
}

func convertITypeToParams(data reflect.Value) Params {
	var params = Params{}

	oType := data.Type()
	for i := 0; i < data.NumField(); i++ {
		el := data.Field(i)
		fieldName := oType.Field(i).Name
		switch el.Type() {
		case reflect.TypeOf(Integer{}):
			var p = Param{}
			p.Name = convertName(fieldName)
			p.Value = el.Interface().(Integer).value
			if p.Value == 0 && !el.Interface().(Integer).write {
				continue
			}
			params = append(params, &p)
		case reflect.TypeOf(String{}):
			var p = Param{}
			p.Name = convertName(fieldName)
			p.Value = el.Interface().(String).value
			if p.Value == "" && !el.Interface().(String).write {
				continue
			}
			params = append(params, &p)
		case reflect.TypeOf(NumberedInteger{}):
			w := 1
			for m, _ := range el.Interface().(NumberedInteger).values {
				var p = Param{}
				p.Name = convertName(fieldName, strconv.Itoa(w))
				p.Value = m
				w += 1
				params = append(params, &p)
			}
		case reflect.TypeOf(NumberedString{}):
			w := 1
			for m, _ := range el.Interface().(NumberedString).values {
				var p = Param{}
				p.Name = convertName(fieldName, strconv.Itoa(w))
				p.Value = m
				w += 1
				params = append(params, &p)
			}
		case reflect.TypeOf(Dict{}):
			for m, _ := range el.Interface().(Dict).values {
				var p = Param{}
				p.Name = convertName(fieldName)
				p.Value = m
				params = append(params, &p)
			}
		case reflect.TypeOf(Array{}):
		default:
			continue
		}
	}

	return params
}
