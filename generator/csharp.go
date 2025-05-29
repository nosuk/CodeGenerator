package generator

import (
	"fmt"
	"strings"

	"github.com/nosuk/CodeGenerator/models"
)

func toCamelCase(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// C# 타입 변환 (배열은 List<>)
func csharpType(field models.Field) string {
	if field.IsArray {
		return fmt.Sprintf("List<%s>", field.Type)
	}
	return field.Type
}

// 배열 타입에서 아이템명 추출 (ex: List<Role> → Role)
func arrayItemType(field models.Field) string {
	t := csharpType(field)
	if strings.HasPrefix(t, "List<") && strings.HasSuffix(t, ">") {
		return t[5 : len(t)-1]
	}
	return t
}

// C# 4.7.2 스타일 코드 생성기 (JSON/XML 동시 지원, 배열/단일 어트리뷰트 자동 분기)
func GenerateCSharpCode(field models.Field, rootClassName string, outputKinds ...OutputKind) string {
	var sb strings.Builder

	sb.WriteString("using System;\nusing System.IO;\nusing System.Collections.Generic;\nusing Newtonsoft.Json;\nusing System.Xml.Serialization;\n\n")

	// 하위 클래스(루트 제외) 정의
	writeCSharpClassTree(field, rootClassName, &sb)

	// 루트 모델 클래스
	sb.WriteString(fmt.Sprintf("[XmlRoot(ElementName=\"%s\")]\n", field.Name))
	sb.WriteString(fmt.Sprintf("public class %s\n{\n", field.Name))
	for _, child := range field.Children {
		if child.IsArray {
			// 배열/리스트: [JsonProperty], [XmlArray], [XmlArrayItem]
			sb.WriteString(fmt.Sprintf("    [JsonProperty(\"%s\")]\n", toCamelCase(child.Name)))
			sb.WriteString(fmt.Sprintf("    [XmlArray(\"%s\")]\n", child.Name))
			sb.WriteString(fmt.Sprintf("    [XmlArrayItem(\"%s\")]\n", arrayItemType(child)))
			sb.WriteString(fmt.Sprintf("    public %s %s { get; set; }\n", csharpType(child), child.Name))
		} else {
			// 단일값: [JsonProperty], [XmlElement]
			sb.WriteString(fmt.Sprintf("    [JsonProperty(\"%s\")]\n", toCamelCase(child.Name)))
			sb.WriteString(fmt.Sprintf("    [XmlElement(\"%s\")]\n", child.Name))
			sb.WriteString(fmt.Sprintf("    public %s %s { get; set; }\n", csharpType(child), child.Name))
		}
	}
	sb.WriteString("}\n\n")

	// IO static class (.NET 4.7.2)
	sb.WriteString(fmt.Sprintf("public static class %sIO\n{\n", rootClassName))

	// JSON 함수
	if HasKind(outputKinds, OutputJSON) {
		sb.WriteString(fmt.Sprintf("    // JSON 입출력\n"))
		sb.WriteString(fmt.Sprintf("    public static %s LoadFromJsonFile(string path)\n", rootClassName))
		sb.WriteString(fmt.Sprintf("    {\n        var json = File.ReadAllText(path);\n        return JsonConvert.DeserializeObject<%s>(json);\n    }\n\n", rootClassName))
		sb.WriteString(fmt.Sprintf("    public static void SaveToJsonFile(string path, %s data)\n", rootClassName))
		sb.WriteString(fmt.Sprintf("    {\n        var json = JsonConvert.SerializeObject(data);\n        File.WriteAllText(path, json);\n    }\n\n"))
		sb.WriteString(fmt.Sprintf("    public static string MarshalJson(%s data)\n", rootClassName))
		sb.WriteString("    {\n        return JsonConvert.SerializeObject(data);\n    }\n\n")
		sb.WriteString(fmt.Sprintf("    public static %s UnmarshalJson(string json)\n", rootClassName))
		sb.WriteString(fmt.Sprintf("    {\n        return JsonConvert.DeserializeObject<%s>(json);\n    }\n\n", rootClassName))
	}

	// XML 함수
	if HasKind(outputKinds, OutputXML) {
		sb.WriteString("    // XML 입출력\n")
		sb.WriteString(fmt.Sprintf(
			"    public static %s LoadFromXmlFile(string path)\n"+
				"    {\n        using (var stream = File.OpenRead(path))\n        {\n            var serializer = new XmlSerializer(typeof(%s));\n            return (%s)serializer.Deserialize(stream);\n        }\n    }\n\n",
			rootClassName, rootClassName, rootClassName,
		))
		sb.WriteString(fmt.Sprintf(
			"    public static void SaveToXmlFile(string path, %s data)\n"+
				"    {\n        using (var stream = File.Create(path))\n        {\n            var serializer = new XmlSerializer(typeof(%s));\n            serializer.Serialize(stream, data);\n        }\n    }\n\n",
			rootClassName, rootClassName,
		))
		sb.WriteString(fmt.Sprintf(
			"    public static string MarshalXml(%s data)\n"+
				"    {\n        using (var ms = new MemoryStream())\n        {\n            var serializer = new XmlSerializer(typeof(%s));\n            serializer.Serialize(ms, data);\n            ms.Position = 0;\n            using (var reader = new StreamReader(ms))\n            {\n                return reader.ReadToEnd();\n            }\n        }\n    }\n\n",
			rootClassName, rootClassName,
		))
		sb.WriteString(fmt.Sprintf(
			"    public static %s UnmarshalXml(string xml)\n"+
				"    {\n        var bytes = System.Text.Encoding.UTF8.GetBytes(xml);\n        using (var ms = new MemoryStream(bytes))\n        {\n            var serializer = new XmlSerializer(typeof(%s));\n            return (%s)serializer.Deserialize(ms);\n        }\n    }\n",
			rootClassName, rootClassName, rootClassName,
		))
	}

	sb.WriteString("}\n\n")

	return sb.String()
}

// 하위 클래스(루트 제외)도 배열/단일 분기 적용해서 생성
func writeCSharpClassTree(field models.Field, rootClassName string, sb *strings.Builder) {
	for _, child := range field.Children {
		if child.IsComplex {
			writeCSharpClassTree(child, rootClassName, sb)
			if child.Name != rootClassName {
				sb.WriteString(fmt.Sprintf("[XmlType(TypeName=\"%s\")]\n", child.Name))
				sb.WriteString(fmt.Sprintf("public class %s\n{\n", child.Name))
				for _, grandChild := range child.Children {
					if grandChild.IsArray {
						sb.WriteString(fmt.Sprintf("    [JsonProperty(\"%s\")]\n", toCamelCase(grandChild.Name)))
						sb.WriteString(fmt.Sprintf("    [XmlArray(\"%s\")]\n", grandChild.Name))
						sb.WriteString(fmt.Sprintf("    [XmlArrayItem(\"%s\")]\n", arrayItemType(grandChild)))
						sb.WriteString(fmt.Sprintf("    public %s %s { get; set; }\n", csharpType(grandChild), grandChild.Name))
					} else {
						sb.WriteString(fmt.Sprintf("    [JsonProperty(\"%s\")]\n", toCamelCase(grandChild.Name)))
						sb.WriteString(fmt.Sprintf("    [XmlElement(\"%s\")]\n", grandChild.Name))
						sb.WriteString(fmt.Sprintf("    public %s %s { get; set; }\n", csharpType(grandChild), grandChild.Name))
					}
				}
				sb.WriteString("}\n\n")
			}
		}
	}
}

// OutputKind, HasKind 등 공통 유틸은 common.go에서 제공
