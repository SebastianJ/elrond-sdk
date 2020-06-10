package api

import "regexp"

var (
	urlValidationRegexp = regexp.MustCompile(`^http(s)?:\/\/wallet\-api\.elrond\.com`)
)

func isWalletAPIURL(url string) bool {
	matches := urlValidationRegexp.FindAllStringSubmatch(url, -1)
	if len(matches) > 0 {
		return true
	}

	return false
}
