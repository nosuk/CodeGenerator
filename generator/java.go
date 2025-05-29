package generator

import (
	"fmt"
	"strings"

	"github.com/nosuk/CodeGenerator/models"
)

// Java 타입 변환 (배열이면 List<타입>)
func javaType(field models.Field) string {
	if field.IsArray {
		return fmt.Sprintf("List<%s>", field.Type)
	}
	return field.Type
}

func GenerateJavaCode(field models.Field, rootClassName string, outputKinds ...OutputKind) string {
	var sb strings.Builder

	// import 구문 (Jackson + JAXB + Java 표준)
	sb.WriteString("import com.fasterxml.jackson.annotation.*;\n")
	sb.WriteString("import com.fasterxml.jackson.databind.ObjectMapper;\n")
	sb.WriteString("import javax.xml.bind.*;\n")
	sb.WriteString("import javax.xml.bind.annotation.*;\n")
	sb.WriteString("import java.io.*;\n")
	sb.WriteString("import java.nio.file.*;\n")
	sb.WriteString("import java.util.*;\n\n")

	// 하위 클래스 (루트 제외)
	writeJavaClassTree(field, rootClassName, &sb)

	// 루트 클래스 정의
	sb.WriteString(fmt.Sprintf("@XmlRootElement(name=\"%s\")\n", field.Name))
	sb.WriteString("@XmlAccessorType(XmlAccessType.FIELD)\n")
	sb.WriteString("@JsonIgnoreProperties(ignoreUnknown=true)\n")
	sb.WriteString(fmt.Sprintf("public class %s {\n", field.Name))
	for _, child := range field.Children {
		if child.IsArray {
			sb.WriteString(fmt.Sprintf("    @XmlElementWrapper(name=\"%s\")\n", child.Name))
			sb.WriteString(fmt.Sprintf("    @XmlElement(name=\"%s\")\n", arrayItemType(child)))
			sb.WriteString(fmt.Sprintf("    @JsonProperty(\"%s\")\n", toCamelCase(child.Name)))
			sb.WriteString(fmt.Sprintf("    public %s %s;\n", javaType(child), child.Name))
		} else {
			sb.WriteString(fmt.Sprintf("    @XmlElement(name=\"%s\")\n", child.Name))
			sb.WriteString(fmt.Sprintf("    @JsonProperty(\"%s\")\n", toCamelCase(child.Name)))
			sb.WriteString(fmt.Sprintf("    public %s %s;\n", javaType(child), child.Name))
		}
	}
	sb.WriteString(fmt.Sprintf("\n    public %s() {}\n", field.Name))
	sb.WriteString("}\n\n")

	// IO 유틸 클래스 (Jackson + JAXB)
	sb.WriteString(fmt.Sprintf("class %sIO {\n", rootClassName))

	// JSON
	sb.WriteString(fmt.Sprintf("    public static %s loadFromJsonFile(String path) throws IOException {\n", rootClassName))
	sb.WriteString("        ObjectMapper om = new ObjectMapper();\n")
	sb.WriteString(fmt.Sprintf("        return om.readValue(Files.readAllBytes(Paths.get(path)), %s.class);\n", rootClassName))
	sb.WriteString("    }\n\n")
	sb.WriteString(fmt.Sprintf("    public static void saveToJsonFile(String path, %s data) throws IOException {\n", rootClassName))
	sb.WriteString("        ObjectMapper om = new ObjectMapper();\n")
	sb.WriteString("        om.writerWithDefaultPrettyPrinter().writeValue(new File(path), data);\n")
	sb.WriteString("    }\n\n")

	// XML
	sb.WriteString(fmt.Sprintf("    public static %s loadFromXmlFile(String path) throws Exception {\n", rootClassName))
	sb.WriteString(fmt.Sprintf("        JAXBContext ctx = JAXBContext.newInstance(%s.class);\n", rootClassName))
	sb.WriteString(fmt.Sprintf("        return (%s) ctx.createUnmarshaller().unmarshal(new File(path));\n", rootClassName))
	sb.WriteString("    }\n\n")
	sb.WriteString(fmt.Sprintf("    public static void saveToXmlFile(String path, %s data) throws Exception {\n", rootClassName))
	sb.WriteString(fmt.Sprintf("        JAXBContext ctx = JAXBContext.newInstance(%s.class);\n", rootClassName))
	sb.WriteString("        ctx.createMarshaller().marshal(data, new File(path));\n")
	sb.WriteString("    }\n")
	sb.WriteString("}\n\n")

	return sb.String()
}

// 하위 클래스도 JSON+XML 어노테이션 포함
func writeJavaClassTree(field models.Field, rootClassName string, sb *strings.Builder) {
	for _, child := range field.Children {
		if child.IsComplex {
			writeJavaClassTree(child, rootClassName, sb)
			if child.Name != rootClassName {
				sb.WriteString(fmt.Sprintf("@XmlType(name=\"%s\")\n", child.Name))
				sb.WriteString("@XmlAccessorType(XmlAccessType.FIELD)\n")
				sb.WriteString("@JsonIgnoreProperties(ignoreUnknown=true)\n")
				sb.WriteString(fmt.Sprintf("public class %s {\n", child.Name))
				for _, grandChild := range child.Children {
					if grandChild.IsArray {
						sb.WriteString(fmt.Sprintf("    @XmlElementWrapper(name=\"%s\")\n", grandChild.Name))
						sb.WriteString(fmt.Sprintf("    @XmlElement(name=\"%s\")\n", arrayItemType(grandChild)))
						sb.WriteString(fmt.Sprintf("    @JsonProperty(\"%s\")\n", toCamelCase(grandChild.Name)))
						sb.WriteString(fmt.Sprintf("    public %s %s;\n", javaType(grandChild), grandChild.Name))
					} else {
						sb.WriteString(fmt.Sprintf("    @XmlElement(name=\"%s\")\n", grandChild.Name))
						sb.WriteString(fmt.Sprintf("    @JsonProperty(\"%s\")\n", toCamelCase(grandChild.Name)))
						sb.WriteString(fmt.Sprintf("    public %s %s;\n", javaType(grandChild), grandChild.Name))
					}
				}
				sb.WriteString(fmt.Sprintf("\n    public %s() {}\n", child.Name))
				sb.WriteString("}\n\n")
			}
		}
	}
}
