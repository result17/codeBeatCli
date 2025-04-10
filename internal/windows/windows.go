package windows

import (
	"regexp"
	"strings"
)

var (
	backslashReplaceRegex = regexp.MustCompile(`[\\/]+`)
	windowsDriveRegex     = regexp.MustCompile("^[a-z]:/")
)

// FormatFilePath formats a windows filepath by converting backslash
func FormatFilePath(fp string) string {
	if windowsDriveRegex.MatchString(fp) {
		fp = strings.ToUpper(fp[:1]) + fp[1:]
	}
	return backslashReplaceRegex.ReplaceAllString(fp, "/")
}
