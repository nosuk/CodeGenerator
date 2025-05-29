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
	// 1️⃣ CLI 인자 처리
	inputPath := flag.String("input", "", "입력 JSON 파일 경로 (예: sample.json)")
	lang := flag.String("lang", "", "타겟 언어 (csharp,go,python 여러개 쉼표 구분)")
	flag.Parse()

	if *inputPath == "" {
		fmt.Println("❗ 입력 파일 경로를 -input 으로 지정해 주세요")
		os.Exit(1)
	}

	// 2️⃣ 입력 파일명에서 기본 정보 추출
	base := filepath.Base(*inputPath)                    // 예: sample.json
	name := strings.TrimSuffix(base, filepath.Ext(base)) // 예: sample
	rootClassName := toExported(name)                    // Sample, Go에서는 struct명, Python에서는 클래스명
	dirName := name                                      // sample

	// 3️⃣ 입력 JSON 파싱 → Field 트리 생성
	data, err := ioutil.ReadFile(*inputPath)
	if err != nil {
		fmt.Println("❗ 파일 읽기 오류:", err)
		os.Exit(1)
	}
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		fmt.Println("❗ JSON 파싱 오류:", err)
		os.Exit(1)
	}
	field := models.ParseJSONToFields(raw, rootClassName)

	// 4️⃣ 언어별 코드 생성 (lang 미입력 시 전체 언어)
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

// 언어별 코드 생성/저장 함수
func generateCodeForLang(lang string, field models.Field, rootClassName, dirName, baseName string) {
	var code string
	var ext string

	switch lang {
	case "csharp":
		code = generator.GenerateCSharpCode(field, rootClassName)
		ext = ".cs"
	case "go":
		code = generator.GenerateGoCode(field, rootClassName)
		ext = ".go"
	case "python":
		code = generator.GeneratePythonCode(field, rootClassName)
		ext = ".py"
	default:
		fmt.Printf("⚠️ 지원하지 않는 언어: %s\n", lang)
		return
	}

	// ./sample/csharp/sample.cs  등 경로 만들기
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

// 첫 글자 대문자로 (Go/C#/Python 클래스 네이밍)
func toExported(name string) string {
	if name == "" {
		return ""
	}
	return strings.ToUpper(name[:1]) + name[1:]
}
