package defines

import "regexp"

var (
	RegexAddPayment = regexp.MustCompile("^\\$*(\\d+) ([a-zA-Z0-9 ]+)(?: \\((.+)\\))*$")
)

