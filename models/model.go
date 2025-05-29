package models

import (
	"encoding/xml"
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

// JSON → Field 트리 (재귀)
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

// XML → Field 트리 (간단 샘플)
// 실제로는 xml.Decoder로 Element별 재귀 파싱/배열/속성 처리 추가 필요
func ParseXMLToFields(data []byte, name string) Field {
	var m map[string]interface{}
	if err := xml.Unmarshal(data, (*map[string]interface{})(&m)); err == nil && len(m) > 0 {
		return ParseJSONToFields(m, name)
	}
	// 실제 구현은 xml.Decoder로 태그 구조 → map 변환 로직 필요
	return Field{Name: ToExported(name), Type: "object"}
}

// 첫글자 대문자 (Go/C#/Python 네이밍)
func ToExported(name string) string {
	if name == "" {
		return ""
	}
	return strings.ToUpper(name[:1]) + name[1:]
}
