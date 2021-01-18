package expr

import "strings"

func Path(path string) []string {
	return strings.Split(path, ".")
}
