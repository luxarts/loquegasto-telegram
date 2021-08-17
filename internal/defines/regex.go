package defines

import "regexp"

var (
	RegexAddPayment = regexp.MustCompile("^\\$*(\\d+(?:(?:\\.|,)\\d+)*) (.+[^()])(?: \\((.+[^()])\\))*$")
)
