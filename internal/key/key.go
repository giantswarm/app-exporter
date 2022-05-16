package key

import "strings"

// formatVersion normalizes version representation by removing `v` prefix.
// It matters for customers Catalogs, ACEs and apps created out of them.
func FormatVersion(input string) string {
	return strings.TrimPrefix(input, "v")
}
