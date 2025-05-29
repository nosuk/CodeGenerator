package generator

import (
	"fmt"
	"strings"

	"github.com/nosuk/CodeGenerator/models"
)

// Python 코드 생성 함수
func GeneratePythonCode(field models.Field, rootName string) string {
	var sb strings.Builder

	// import문
	sb.WriteString("import json\n\n")

	// 클래스 정의 (하위 클래스부터)
	writePythonClasses(field, rootName, &sb)

	// 파일 입출력 및 직렬화/역직렬화 함수 (루트만)
	sb.WriteString(fmt.Sprintf("def load_%s_from_file(path):\n", to_snake_case(rootName)))
	sb.WriteString(fmt.Sprintf("    with open(path, 'r', encoding='utf-8') as f:\n        data = json.load(f)\n    return %s.from_dict(data)\n\n", rootName))

	sb.WriteString(fmt.Sprintf("def save_%s_to_file(path, obj):\n", to_snake_case(rootName)))
	sb.WriteString("    with open(path, 'w', encoding='utf-8') as f:\n        json.dump(obj.to_dict(), f, ensure_ascii=False, indent=2)\n\n")

	sb.WriteString(fmt.Sprintf("def marshal_%s(obj):\n    return json.dumps(obj.to_dict(), ensure_ascii=False)\n\n", to_snake_case(rootName)))
	sb.WriteString(fmt.Sprintf("def unmarshal_%s(json_str):\n    data = json.loads(json_str)\n    return %s.from_dict(data)\n\n", to_snake_case(rootName), rootName))

	return sb.String()
}

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

// 클래스 본체
func writePythonClass(field models.Field, sb *strings.Builder) {
	// __init__
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

	// from_dict
	sb.WriteString("    @staticmethod\n")
	sb.WriteString("    def from_dict(obj):\n")
	sb.WriteString("        if obj is None: return None\n")
	sb.WriteString(fmt.Sprintf("        return %s(\n", field.Name))
	for i, c := range field.Children {
		// 중첩 구조는 재귀 호출
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

	// to_dict
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

// snake_case 변환
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

// 마지막에 콤마 찍을지 여부
func if_comma(i int, l []models.Field) string {
	if i != len(l)-1 {
		return ","
	}
	return ""
}
