package patch

import (
	"sort"
	"strconv"
	"strings"
)

// Tags represents a simple key-value type for struct tags.
type Tags map[string]string

func (tags Tags) String() string {
	// Sort tag names for deterministic output
	names := make([]string, 0, len(tags))
	for name := range tags {
		names = append(names, name)
	}
	sort.Strings(names)
	var b strings.Builder
	for i, name := range names {
		if i > 0 {
			b.WriteRune(' ')
		}
		b.WriteString(name)
		b.WriteRune(':')
		b.WriteString(strconv.Quote(tags[name]))
	}
	return b.String()
}
