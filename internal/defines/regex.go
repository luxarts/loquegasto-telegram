package defines

import "regexp"

var (
	RegexTransaction = regexp.MustCompile("^\\$*(\\-*\\d+(?:(?:\\.|,)\\d+)*) (.+[^()])(?: \\((.+[^()])\\))*$")
)
