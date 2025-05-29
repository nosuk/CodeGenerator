package generator

import (
	"fmt"
	"strings"

	"github.com/nosuk/CodeGenerator/models"
)

// Go 타입 변환: 배열이면 []타입, 아니면 타입명
func goType(field models.Field) string {
	if field.IsArray {
		return "[]" + field.Type
	}
	return field.Type
}

// Go 코드 생성기 (JSON/XML 동시 지원)
func GenerateGoCode(field models.Field, rootName string, outputKinds ...OutputKind) string {
	var sb strings.Builder
	sb.WriteString("package main\n\n")
	sb.WriteString("import (\n\t\"encoding/json\"\n\t\"encoding/xml\"\n\t\"os\"\n\t\"io/ioutil\"\n)\n\n")

	// struct 정의 (하위 struct 먼저)
	writeGoStructs(field, rootName, &sb)

	// JSON 입출력
	if HasKind(outputKinds, OutputJSON) {
		sb.WriteString(fmt.Sprintf("// 파일에서 JSON 읽기\nfunc Load%sFromJSONFile(path string) (%s, error) {\n", rootName, rootName))
		sb.WriteString(fmt.Sprintf("    var v %s\n", rootName))
		sb.WriteString("    data, err := ioutil.ReadFile(path)\n    if err != nil { return v, err }\n")
		sb.WriteString("    err = json.Unmarshal(data, &v)\n    return v, err\n}\n\n")

		sb.WriteString(fmt.Sprintf("// JSON 파일로 저장\nfunc Save%sToJSONFile(path string, v %s) error {\n", rootName, rootName))
		sb.WriteString("    data, err := json.MarshalIndent(v, \"\", \"  \")\n    if err != nil { return err }\n")
		sb.WriteString("    return ioutil.WriteFile(path, data, 0644)\n}\n\n")
	}

	// XML 입출력
	if HasKind(outputKinds, OutputXML) {
		sb.WriteString(fmt.Sprintf("// 파일에서 XML 읽기\nfunc Load%sFromXMLFile(path string) (%s, error) {\n", rootName, rootName))
		sb.WriteString(fmt.Sprintf("    var v %s\n", rootName))
		sb.WriteString("    data, err := ioutil.ReadFile(path)\n    if err != nil { return v, err }\n")
		sb.WriteString("    err = xml.Unmarshal(data, &v)\n    return v, err\n}\n\n")

		sb.WriteString(fmt.Sprintf("// XML 파일로 저장\nfunc Save%sToXMLFile(path string, v %s) error {\n", rootName, rootName))
		sb.WriteString("    data, err := xml.MarshalIndent(v, \"\", \"  \")\n    if err != nil { return err }\n")
		sb.WriteString("    return ioutil.WriteFile(path, data, 0644)\n}\n\n")
	}
	return sb.String()
}

// 하위 struct(루트 제외) 재귀 생성
func writeGoStructs(field models.Field, rootName string, sb *strings.Builder) {
	for _, child := range field.Children {
		if child.IsComplex {
			writeGoStructs(child, rootName, sb)
			if child.Name != rootName {
				sb.WriteString(fmt.Sprintf("type %s struct {\n", child.Name))
				for _, gc := range child.Children {
					sb.WriteString(fmt.Sprintf("    %s %s `json:\"%s\" xml:\"%s\"`\n", gc.Name, goType(gc), gc.Name, gc.Name))
				}
				sb.WriteString("}\n\n")
			}
		}
	}
	// 마지막에 루트 struct 생성
	if field.Name == rootName {
		sb.WriteString(fmt.Sprintf("type %s struct {\n", field.Name))
		for _, child := range field.Children {
			sb.WriteString(fmt.Sprintf("    %s %s `json:\"%s\" xml:\"%s\"`\n", child.Name, goType(child), child.Name, child.Name))
		}
		sb.WriteString("}\n\n")
	}
}

// OutputKind 체크는 generator/common.go에서 제공 (import해서 사용)
