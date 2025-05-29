package generator

import (
	"strings"

	"github.com/nosuk/CodeGenerator/models"
)

type OutputKind string

const (
	OutputJSON OutputKind = "json"
	OutputXML  OutputKind = "xml"
)

// outputKinds(옵션)에 포함 여부 체크
func HasKind(kinds []OutputKind, kind OutputKind) bool {
	if len(kinds) == 0 {
		return true // 옵션 비었으면 모두 지원
	}
	for _, k := range kinds {
		if k == kind {
			return true
		}
	}
	return false
}

func toCamelCase(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func arrayItemType(field models.Field) string {
	t := csharpType(field)
	if strings.HasPrefix(t, "List<") && strings.HasSuffix(t, ">") {
		return t[5 : len(t)-1]
	}
	return t
}
