package defines

import "regexp"

var (
	RegexTransaction  = regexp.MustCompile("^\\$*(\\-*\\d+(?:(?:\\.|,)\\d+)*) (.+[^()])(?: \\((.+[^()])\\))*$")
	RegexCreateWallet = regexp.MustCompile("^(.+)\\s\\$*(\\-*\\d+(?:(?:\\.|,)\\d+)*)$")
)
