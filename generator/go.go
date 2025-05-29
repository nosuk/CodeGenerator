package generator

import (
	"fmt"
	"strings"

	"github.com/nosuk/CodeGenerator/models"
)

// Go 코드 생성 함수
func GenerateGoCode(field models.Field, rootName string) string {
	var sb strings.Builder

	// 패키지명
	sb.WriteString("package main\n\n")
	sb.WriteString("import (\n\t\"encoding/json\"\n\t\"os\"\n\t\"io/ioutil\"\n)\n\n")

	// struct 정의 (하위 struct 먼저)
	writeGoStructs(field, rootName, &sb)

	// 파일 입출력 + Marshal/Unmarshal 함수(루트 struct에만)
	sb.WriteString(fmt.Sprintf("// 파일에서 읽어서 구조체로 파싱\nfunc Load%sFromFile(path string) (%s, error) {\n", rootName, rootName))
	sb.WriteString(fmt.Sprintf("    var v %s\n", rootName))
	sb.WriteString("    data, err := ioutil.ReadFile(path)\n    if err != nil {\n        return v, err\n    }\n")
	sb.WriteString("    err = json.Unmarshal(data, &v)\n    return v, err\n}\n\n")

	sb.WriteString(fmt.Sprintf("// 구조체를 파일로 저장\nfunc Save%sToFile(path string, v %s) error {\n", rootName, rootName))
	sb.WriteString("    data, err := json.MarshalIndent(v, \"\", \"  \")\n    if err != nil {\n        return err\n    }\n")
	sb.WriteString("    return ioutil.WriteFile(path, data, 0644)\n}\n\n")

	sb.WriteString(fmt.Sprintf("// Marshal 구조체→[]byte\nfunc Marshal%s(v %s) ([]byte, error) {\n    return json.Marshal(v)\n}\n\n", rootName, rootName))
	sb.WriteString(fmt.Sprintf("// Unmarshal []byte→구조체\nfunc Unmarshal%s(data []byte) (%s, error) {\n", rootName, rootName))
	sb.WriteString(fmt.Sprintf("    var v %s\n", rootName))
	sb.WriteString("    err := json.Unmarshal(data, &v)\n    return v, err\n}\n\n")
	return sb.String()
}

func writeGoStructs(field models.Field, rootName string, sb *strings.Builder) {
	for _, child := range field.Children {
		if child.IsComplex {
			writeGoStructs(child, rootName, sb)
			if child.Name != rootName {
				sb.WriteString(fmt.Sprintf("type %s struct {\n", child.Name))
				for _, gc := range child.Children {
					sb.WriteString(fmt.Sprintf("    %s %s `json:\"%s\"`\n", gc.Name, gc.Type, gc.Name))
				}
				sb.WriteString("}\n\n")
			}
		}
	}
	// 마지막에 루트 struct 생성
	if field.Name == rootName {
		sb.WriteString(fmt.Sprintf("type %s struct {\n", field.Name))
		for _, child := range field.Children {
			sb.WriteString(fmt.Sprintf("    %s %s `json:\"%s\"`\n", child.Name, child.Type, child.Name))
		}
		sb.WriteString("}\n\n")
	}
}
