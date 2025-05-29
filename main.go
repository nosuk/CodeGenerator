package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/nosuk/CodeGenerator/generator"
	"github.com/nosuk/CodeGenerator/models"
)

func main() {
	inputPath := flag.String("input", "", "입력 파일 경로 (예: sample.json, sample.xml)")
	lang := flag.String("lang", "", "타겟 언어 (csharp,go,python 여러개 쉼표 구분)")
	flag.Parse()

	if *inputPath == "" {
		fmt.Println("❗ 입력 파일 경로를 -input 으로 지정해 주세요")
		os.Exit(1)
	}

	base := filepath.Base(*inputPath)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	rootClassName := models.ToExported(name)
	dirName := name

	data, err := ioutil.ReadFile(*inputPath)
	if err != nil {
		fmt.Println("❗ 파일 읽기 오류:", err)
		os.Exit(1)
	}

	// 1️⃣ 확장자 감지로 JSON/XML 파싱 분기
	ext := strings.ToLower(filepath.Ext(*inputPath))
	var field models.Field
	if ext == ".json" {
		var raw interface{}
		if err := json.Unmarshal(data, &raw); err != nil {
			fmt.Println("❗ JSON 파싱 오류:", err)
			os.Exit(1)
		}
		field = models.ParseJSONToFields(raw, rootClassName)
	} else if ext == ".xml" {
		field = models.ParseXMLToFields(data, rootClassName)
	} else {
		fmt.Println("지원하지 않는 입력 파일 형식입니다.")
		os.Exit(1)
	}

	langs := []string{"csharp", "go", "python"}

	if *lang == "" {
		for _, l := range langs {
			generateCodeForLang(l, field, rootClassName, dirName, name)
		}
	} else {
		for _, l := range strings.Split(*lang, ",") {
			generateCodeForLang(strings.ToLower(strings.TrimSpace(l)), field, rootClassName, dirName, name)
		}
	}
}

// generator/아래 OutputKind와 일치해야 함!
type OutputKind string

const (
	OutputJSON OutputKind = "json"
	OutputXML  OutputKind = "xml"
)

// 언어별 코드 생성/저장 함수
func generateCodeForLang(lang string, field models.Field, rootClassName, dirName, baseName string) {
	var code string
	var ext string
	kinds := []generator.OutputKind{generator.OutputJSON, generator.OutputXML} // 항상 둘 다 지원

	switch lang {
	case "csharp":
		code = generator.GenerateCSharpCode(field, rootClassName, kinds...)
		ext = ".cs"
	case "go":
		code = generator.GenerateGoCode(field, rootClassName, kinds...)
		ext = ".go"
	case "python":
		code = generator.GeneratePythonCode(field, rootClassName, kinds...)
		ext = ".py"
	default:
		fmt.Printf("⚠️ 지원하지 않는 언어: %s\n", lang)
		return
	}

	targetDir := filepath.Join(".", dirName, lang)
	targetPath := filepath.Join(targetDir, baseName+ext)

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		fmt.Println("❗ 디렉토리 생성 오류:", err)
		return
	}
	if err := ioutil.WriteFile(targetPath, []byte(code), 0644); err != nil {
		fmt.Println("❗ 파일 저장 오류:", err)
		return
	}

	fmt.Printf("✅ %s 코드 생성 완료: %s\n", lang, targetPath)
}
