# 다중 언어 코드 생성기 (Multi-Language Code Generator)

## 1. 개요 (Overview)
- **프로젝트 명**: 다중 언어 코드 생성기
- **주요 내용**:  
  JSON 및 XML 텍스트 입력을 받아, 여러 프로그래밍 언어(Java, Python, Go, C#, C++, JavaScript)별 데이터 구조 정의, 마샬/언마샬, 파일 입출력 코드 자동 생성  
- **주요 기능**:  
  - JSON/XML 파싱 및 공통 데이터 모델 생성  
  - 언어별 데이터 구조 코드 생성 (클래스, struct 등)  
  - 마샬링/언마샬링 코드 생성  
  - 파일 I/O 코드 생성  
  - 생성된 파일을 프로젝트에 바로 등록 가능하도록 지원  

## 2. 목적 (Purpose)
- **자동화 목표**:  
  반복적이고 수작업이 많은 데이터 모델 정의 및 입출력 코드 작성을 자동화하여 개발 생산성 극대화  
- **효과 및 기대효과**:  
  - 빠르고 일관된 데이터 구조 생성  
  - 다양한 언어 지원으로 멀티 플랫폼 개발 대응  
  - 인적 오류 감소 및 유지보수 효율성 증가  
- **대상 사용자**:  
  - 개발자 (특히 Golang, C# 등 다중 언어 환경)  
  - 시스템 통합 및 자동화 전문가  

## 3. 설계 내용 (Design Overview)

### 3-1. 입력 처리
- JSON 및 XML 텍스트 파일 또는 문자열 입력 지원  
- 입력 데이터 파싱 후 공통 데이터 모델(트리 구조) 생성  

### 3-2. 공통 데이터 모델
- 필드명, 타입, 중첩 구조(복합 타입), 배열 여부 등 정보 포함  
- 중첩된 구조체/클래스 자동 인식 및 재귀 처리  

### 3-3. 코드 생성 모듈
- 타겟 언어별 템플릿 보유  
- 데이터 구조 정의 + 마샬/언마샬 + 파일 입출력 함수 포함  
- Go 템플릿 기반 확장성 있는 설계  

### 3-4. 출력 및 파일 관리
- 각 언어별 소스 코드 파일 생성  
- 프로젝트 내 쉽게 포함할 수 있도록 파일 네이밍 및 경로 관리  
- CLI 옵션 또는 GUI 인터페이스 제공 계획  

## 4. 상세 설계 (Detailed Design)

### 4-1. 입력 파싱
- JSON: `encoding/json` 재귀 파싱 → map[string]interface{} → Field 구조체 변환  
- XML: `encoding/xml` 재귀 파싱 → 비슷한 공통 데이터 모델 변환  

### 4-2. 데이터 모델 구조 (Field struct)
- Name: 필드명 (대문자 변환 포함)  
- Type: 데이터 타입 (언어별 변환 가능)  
- JsonTag / XmlTag: 직렬화 시 필드명  
- Children: 중첩 필드 리스트  
- IsArray, IsComplex: 배열 여부 및 복합 타입 표시  

### 4-3. 언어별 템플릿
- Java: 클래스 + Jackson/Gson 마샬/언마샬 + 파일 입출력  
- Python: 클래스/딕셔너리 + json/xml 표준 라이브러리 + 파일 입출력  
- Go: struct + encoding/json, encoding/xml + 파일 I/O  
- C#: 클래스 + System.Text.Json, XmlSerializer + 파일 I/O  
- C++: struct/class + nlohmann/json 라이브러리 + fstream  
- JavaScript(Node.js): 객체 + JSON/XML 파싱 라이브러리 + fs 모듈  

### 4-4. 확장성 고려사항
- 중첩 구조 및 배열 처리 강화  
- 사용자 정의 타입 및 커스텀 매핑 지원  
- CLI 인자, 파일 입출력 경로 지정 지원  
- 다양한 직렬화 라이브러리 옵션 선택 가능  

## 5. 코드 실행 방식 및 파라미터 설계

### 5-1. 입력 파라미터
- `-input` : JSON 또는 XML 파일 경로 (필수)  
- `-lang` : 생성할 언어 지정 (복수 가능, 예: go,csharp,python)  
  - 입력 없으면 기본값으로 모든 지원 언어 전체 생성  
- `-output` : 생성할 파일명 또는 출력 디렉터리 지정  
  - 파일명 지정 시 해당 파일 1개 생성  
  - 지정 없으면 하위 폴더를 만들고 각 언어별로 전체 파일 생성 (예: ./output/go/, ./output/csharp/)  

### 5-2. 실행 예시
```bash
# JSON 파일을 입력받아 Go와 C# 코드 생성, 출력 폴더 지정
./codegen -input sample.json -lang go,csharp -output ./generated_codes

# XML 파일 입력 후 모든 지원 언어 전체 코드 생성, 기본 폴더 구조로 출력
./codegen -input sample.xml
