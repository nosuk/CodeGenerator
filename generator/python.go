package generator

import (
	"fmt"
	"strings"

	"github.com/nosuk/CodeGenerator/models"
)

// Python 코드 생성기 - JSON, XML 지원
func GeneratePythonCode(field models.Field, rootName string, outputKinds ...OutputKind) string {
	var sb strings.Builder

	// import문
	sb.WriteString("import json\n")
	sb.WriteString("import xml.etree.ElementTree as ET\n\n")

	// 클래스 정의 (하위 클래스부터)
	writePythonClasses(field, rootName, &sb)

	// JSON 함수
	if HasKind(outputKinds, OutputJSON) {
		sb.WriteString(fmt.Sprintf("def load_%s_from_json_file(path):\n", to_snake_case(rootName)))
		sb.WriteString(fmt.Sprintf("    with open(path, 'r', encoding='utf-8') as f:\n        data = json.load(f)\n    return %s.from_dict(data)\n\n", rootName))
		sb.WriteString(fmt.Sprintf("def save_%s_to_json_file(path, obj):\n", to_snake_case(rootName)))
		sb.WriteString("    with open(path, 'w', encoding='utf-8') as f:\n        json.dump(obj.to_dict(), f, ensure_ascii=False, indent=2)\n\n")
	}

	// XML 함수(간단 버전: xml.etree.ElementTree 이용)
	if HasKind(outputKinds, OutputXML) {
		sb.WriteString(fmt.Sprintf("# XML 지원은 기본 dict 변환을 가정한 예시, 실전용 구현은 확장 필요\n"))
		sb.WriteString(fmt.Sprintf("def load_%s_from_xml_file(path):\n", to_snake_case(rootName)))
		sb.WriteString("    tree = ET.parse(path)\n    root = tree.getroot()\n    # TODO: ElementTree → dict → 클래스 변환 구현 필요\n\n")
		sb.WriteString(fmt.Sprintf("def save_%s_to_xml_file(path, obj):\n", to_snake_case(rootName)))
		sb.WriteString("    # TODO: 클래스 → dict → ElementTree 변환 구현 필요\n    pass\n\n")
	}

	return sb.String()
}

// ... 이하 writePythonClasses, to_snake_case, HasKind 함수는 동일
func writePythonClasses(field models.Field, rootName string, sb *strings.Builder) {
	for _, child := range field.Children {
		if child.IsComplex {
			writePythonClasses(child, rootName, sb)
			if child.Name != rootName {
				writePythonClass(child, sb)
			}
		}
	}
	if field.Name == rootName {
		writePythonClass(field, sb)
	}
}

func writePythonClass(field models.Field, sb *strings.Builder) {
	sb.WriteString(fmt.Sprintf("class %s:\n", field.Name))
	sb.WriteString("    def __init__(self")
	for _, c := range field.Children {
		sb.WriteString(fmt.Sprintf(", %s=None", to_snake_case(c.Name)))
	}
	sb.WriteString("):\n")
	for _, c := range field.Children {
		sb.WriteString(fmt.Sprintf("        self.%s = %s\n", to_snake_case(c.Name), to_snake_case(c.Name)))
	}
	sb.WriteString("\n")

	sb.WriteString("    @staticmethod\n")
	sb.WriteString("    def from_dict(obj):\n")
	sb.WriteString("        if obj is None: return None\n")
	sb.WriteString(fmt.Sprintf("        return %s(\n", field.Name))
	for i, c := range field.Children {
		if c.IsComplex {
			if c.IsArray {
				sb.WriteString(fmt.Sprintf("            %s=[%s.from_dict(x) for x in obj.get('%s', [])]%s\n", to_snake_case(c.Name), c.Name, c.Name, if_comma(i, field.Children)))
			} else {
				sb.WriteString(fmt.Sprintf("            %s=%s.from_dict(obj.get('%s'))%s\n", to_snake_case(c.Name), c.Name, c.Name, if_comma(i, field.Children)))
			}
		} else {
			sb.WriteString(fmt.Sprintf("            %s=obj.get('%s')%s\n", to_snake_case(c.Name), c.Name, if_comma(i, field.Children)))
		}
	}
	sb.WriteString("        )\n\n")

	sb.WriteString("    def to_dict(self):\n")
	sb.WriteString("        result = {}\n")
	for _, c := range field.Children {
		if c.IsComplex && c.IsArray {
			sb.WriteString(fmt.Sprintf("        result['%s'] = [x.to_dict() for x in self.%s] if self.%s is not None else []\n", c.Name, to_snake_case(c.Name), to_snake_case(c.Name)))
		} else if c.IsComplex {
			sb.WriteString(fmt.Sprintf("        result['%s'] = self.%s.to_dict() if self.%s else None\n", c.Name, to_snake_case(c.Name), to_snake_case(c.Name)))
		} else {
			sb.WriteString(fmt.Sprintf("        result['%s'] = self.%s\n", c.Name, to_snake_case(c.Name)))
		}
	}
	sb.WriteString("        return result\n\n")
}

func to_snake_case(s string) string {
	var out []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			out = append(out, '_')
		}
		out = append(out, r)
	}
	return strings.ToLower(string(out))
}

func if_comma(i int, l []models.Field) string {
	if i != len(l)-1 {
		return ","
	}
	return ""
}
