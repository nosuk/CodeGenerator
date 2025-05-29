package generator

import (
	"fmt"
	"strings"

	"github.com/nosuk/CodeGenerator/models"
)

// C# 코드 생성기: using문, 클래스, IO static class(루트만) 생성
func GenerateCSharpCode(field models.Field, rootClassName string) string {
	var sb strings.Builder

	// using문
	sb.WriteString("using System;\nusing System.IO;\nusing System.Collections.Generic;\nusing Newtonsoft.Json;\n\n")

	// 하위 클래스(루트 제외) 정의
	writeClassTree(field, rootClassName, &sb)

	// 루트 모델(파일명과 동일)
	sb.WriteString(fmt.Sprintf("public class %s\n{\n", field.Name))
	for _, child := range field.Children {
		sb.WriteString(fmt.Sprintf("    public %s %s { get; set; }\n", child.Type, child.Name))
	}
	sb.WriteString("}\n\n")

	// 루트 IO static class만 추가
	sb.WriteString(fmt.Sprintf("public static class %sIO\n{\n", rootClassName))
	sb.WriteString(fmt.Sprintf("    public static %s LoadFromFile(string path)\n", rootClassName))
	sb.WriteString(fmt.Sprintf("        => JsonConvert.DeserializeObject<%s>(File.ReadAllText(path));\n\n", rootClassName))
	sb.WriteString(fmt.Sprintf("    public static %s LoadFromString(string json)\n", rootClassName))
	sb.WriteString(fmt.Sprintf("        => JsonConvert.DeserializeObject<%s>(json);\n\n", rootClassName))
	sb.WriteString(fmt.Sprintf("    public static void SaveToFile(string path, %s data)\n", rootClassName))
	sb.WriteString("        => File.WriteAllText(path, JsonConvert.SerializeObject(data));\n\n")
	sb.WriteString(fmt.Sprintf("    public static string Marshal(%s data)\n", rootClassName))
	sb.WriteString("        => JsonConvert.SerializeObject(data);\n\n")
	sb.WriteString(fmt.Sprintf("    public static %s Unmarshal(string json)\n", rootClassName))
	sb.WriteString(fmt.Sprintf("        => JsonConvert.DeserializeObject<%s>(json);\n", rootClassName))
	sb.WriteString("}\n\n")

	return sb.String()
}

// 하위 클래스 정의(루트 제외, 재귀)
func writeClassTree(field models.Field, rootClassName string, sb *strings.Builder) {
	for _, child := range field.Children {
		if child.IsComplex {
			writeClassTree(child, rootClassName, sb)
			if child.Name != rootClassName {
				sb.WriteString(fmt.Sprintf("public class %s\n{\n", child.Name))
				for _, grandChild := range child.Children {
					sb.WriteString(fmt.Sprintf("    public %s %s { get; set; }\n", grandChild.Type, grandChild.Name))
				}
				sb.WriteString("}\n\n")
			}
		}
	}
}
