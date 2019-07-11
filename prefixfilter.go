package complete

import (
	"strings"
)

// PrefixFilter filters a predictor's options based on the prefix
type PrefixFilter interface {
	FilterPrefix(str, prefix string) bool
}

// PrefixFilteringPredictor is a Predictor that also implements PrefixFilter
type PrefixFilteringPredictor struct {
	Predictor        Predictor
	PrefixFilterFunc func(s, prefix string) bool
}

func (p *PrefixFilteringPredictor) Predict(a Args) []string {
	if p.Predictor == nil {
		return []string{}
	}
	return p.Predictor.Predict(a)
}

func (p *PrefixFilteringPredictor) FilterPrefix(str, prefix string) bool {
	if p.PrefixFilterFunc == nil {
		return defaultPrefixFilter(str, prefix)
	}
	return p.PrefixFilterFunc(str, prefix)
}

// defaultPrefixFilter is the PrefixFilter used when none is set
func defaultPrefixFilter(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// PermissivePrefixFilter always returns true
func PermissivePrefixFilter(_, _ string) bool {
	return true
}

// CaseInsensitivePrefixFilter ignores case differences between the prefix and tested string
func CaseInsensitivePrefixFilter(s, prefix string) bool {
	if len(prefix) > len(s) {
		return false
	}
	return strings.EqualFold(prefix, s[:len(prefix)])
}
