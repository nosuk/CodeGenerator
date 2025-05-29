package models

import (
	"strings"
)

// 데이터 구조 트리
type Field struct {
	Name      string
	Type      string
	Children  []Field
	IsArray   bool
	IsComplex bool
}

// JSON 데이터를 Field 트리로 변환 (재귀)
func ParseJSONToFields(data interface{}, name string) Field {
	switch v := data.(type) {
	case map[string]interface{}:
		children := []Field{}
		for key, value := range v {
			childField := ParseJSONToFields(value, key)
			children = append(children, childField)
		}
		return Field{
			Name:      ToExported(name),
			Type:      ToExported(name),
			Children:  children,
			IsArray:   false,
			IsComplex: true,
		}
	case []interface{}:
		if len(v) > 0 {
			childField := ParseJSONToFields(v[0], name)
			return Field{
				Name:      ToExported(name),
				Type:      "List<" + childField.Type + ">",
				Children:  childField.Children,
				IsArray:   true,
				IsComplex: childField.IsComplex,
			}
		} else {
			return Field{
				Name:      ToExported(name),
				Type:      "List<object>",
				Children:  nil,
				IsArray:   true,
				IsComplex: false,
			}
		}
	case string:
		return Field{Name: ToExported(name), Type: "string"}
	case float64:
		return Field{Name: ToExported(name), Type: "int"}
	case bool:
		return Field{Name: ToExported(name), Type: "bool"}
	default:
		return Field{Name: ToExported(name), Type: "object"}
	}
}

// 필드명 대문자(Exported)로 변환
func ToExported(name string) string {
	if name == "" {
		return ""
	}
	return strings.ToUpper(name[:1]) + name[1:]
}
