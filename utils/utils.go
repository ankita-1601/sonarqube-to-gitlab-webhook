package utils

import "regexp"

//StringInSlice checks if a slice contains a specific string
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// CleanNextLinksHeader func
func CleanNextLinksHeader(links string) string {
	rgx := regexp.MustCompile(`\<(.*?)\>`)
	rs := rgx.FindStringSubmatch(links)
	return rs[1]
}
