package host

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func TestWildcardHost(t *testing.T) {
	host := "*.letsgo.tech"
	okHosts := []string{
		"user.letsgo.tech",
		"blog.letsgo.tech",
	}
	errorHosts := []string{
		"www.google.com",
		"www.baidu.com",
	}
	regexStr := ""

	if strings.Contains(host, "*") {
		regexStr = strings.Replace(host, ".", "\\.", -1)
		regexStr = strings.Replace(regexStr, "*", ".+", -1)
	}

	r, _ := regexp.Compile(fmt.Sprintf("^%s$", regexStr))

	for _, h := range okHosts {
		if !r.MatchString(h) {
			t.Error("okHosts matched error")
		}
	}

	for _, h := range errorHosts {
		if r.MatchString(h) {
			t.Error("errorHosts shouldn't match success")
		}
	}
}
