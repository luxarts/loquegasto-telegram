package defines

import "regexp"

var (
	RegexPayment = regexp.MustCompile("^\\$*(\\d+(?:(?:\\.|,)\\d+)*) (.+[^()])(?: \\((.+[^()])\\))*$")
)
