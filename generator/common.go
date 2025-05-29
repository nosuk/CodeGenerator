package generator

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
