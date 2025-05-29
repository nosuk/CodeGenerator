# CodeGenerator

다중 언어 데이터 모델 코드 자동 생성기  
**입력된 JSON 파일을 바탕으로 C#, Go, Python**  
3개 언어의 데이터 모델 코드, 직렬화/역직렬화, 파일 I/O 함수를  
프로젝트 구조에 맞게 자동으로 생성해줍니다.

---

## ✨ 주요 특징

- **JSON 입력만으로 C#, Go, Python 코드 자동 생성**
- 중첩 구조, 배열 등 복합 타입 완벽 지원 (재귀적 분석)
- 각 언어별 네이밍/관례에 맞는 클래스(struct) 코드, 마샬/언마샬, 파일 입출력 함수 포함
- 결과 파일은 `./입력파일명/언어/입력파일명.확장자` 구조로 자동 저장  
  예:  
  - `./sample/csharp/sample.cs`  
  - `./sample/go/sample.go`  
  - `./sample/python/sample.py`
- 협업 및 확장성을 고려한 패키지/모듈 구조

---

## 🚀 설치 및 빌드

1. Go 1.18+ 설치 ([Go 다운로드](https://go.dev/dl/))
2. 저장소 클론
    ```bash
    git clone https://github.com/nosuk/CodeGenerator.git
    cd CodeGenerator
    ```
3. 빌드
    ```bash
    go build -o codegen main.go
    ```

---

## ⚡ 사용법

### 기본 사용 예시
```bash
./codegen -input sample.json
```
입력 JSON 파일을 바탕으로 C#, Go, Python 코드가 자동 생성됩니다.

### 특정 언어만 생성
```bash
./codegen -input sample.json -lang go
./codegen -input sample.json -lang csharp,python
```
- 쉼표로 여러 언어 지정 가능  
- 지원 언어: `csharp`, `go`, `python`

### 결과 파일 구조
```
./sample/csharp/sample.cs
./sample/go/sample.go
./sample/python/sample.py
```

---

## 🛠️ 프로젝트 구조

- `main.go` – CLI 및 실행 진입점  
- `models/` – Field 구조체, JSON 파싱, 공통 유틸  
- `generator/` – 언어별 코드 생성 모듈  
  - `csharp.go` – C# (Newtonsoft.Json 기반)  
  - `go.go` – Go (encoding/json 사용)  
  - `python.go` – Python (표준 json 모듈 사용)

---

## 📋 샘플 입력/출력

### 입력 JSON 예시
```json
{
  "user": {
    "id": 1,
    "name": "Alice",
    "tags": ["dev", "ops"]
  },
  "roles": ["admin", "user"],
  "isActive": true
}
```

### 출력 파일 예시
- `./sample/csharp/sample.cs`
- `./sample/go/sample.go`
- `./sample/python/sample.py`

→ 각 언어별 데이터 모델, 마샬/언마샬, 파일 IO 함수 자동 생성

---

## 📦 패키지 설치 (필요 시)

C# 코드 사용 시 Newtonsoft.Json 패키지 필요:
```bash
dotnet add package Newtonsoft.Json
```

---

## ✅ TODO

- Java, TypeScript, Kotlin 등 언어 추가 예정
- 네임스페이스, JsonProperty 등 고급 옵션 지원
- 커스텀 타입 매핑 및 유닛테스트 강화
