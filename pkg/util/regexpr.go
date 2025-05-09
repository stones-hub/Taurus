package util

import (
	"regexp"
)

const (
	IPV4 = iota
	IPV6
)

// Match expr 正则， ss 待正则匹配的字符串
func Match(expr string, ss string) []string {
	var matcher = regexp.MustCompile(expr)
	return matcher.FindAllString(ss, -1)
}

// MatchEmail 匹配邮箱
func MatchEmail(email string) []string {
	matcher := regexp.MustCompile(`^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(.[a-zA-Z0-9_-]+)+$`)
	return matcher.FindAllString(email, -1)
}

// MatchPhone 匹配手机号
func MatchPhone(iphone string) []string {
	matcher := regexp.MustCompile(`^1(3|4|5|6|7|8|9)[0-9]{9}$`)
	return matcher.FindAllString(iphone, -1)
}

// MatchDomain 匹配域名
func MatchDomain(domain string) []string {
	matcher := regexp.MustCompile(`((([A-Za-z]{3,9}:(?:\/\/)?)(?:[-;:&=+\$,\w]+@)?[A-Za-z0-9.-]+(:[0-9]+)?|(?:www\.|[-;:&=+\$,\w]+@)[A-Za-z0-9.-]+)((?:\/[+~%\/.\w-_]*)?\??(?:[-+=&;%@.\w_]*)#?(?:[\w]*))?)`)
	return matcher.FindAllString(domain, -1)
}

// MatchIP 匹配IP
func MatchIP(ip string, t int) []string {
	var matcher *regexp.Regexp
	if t == IPV4 {
		matcher = regexp.MustCompile(`^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)`)
	} else {
		matcher = regexp.MustCompile(`^([\da-fA-F]{1,4}:){7}[\da-fA-F]{1,4}|:((:[\da−fA−F]1,4)1,6|:)`)
	}
	return matcher.FindAllString(ip, -1)
}

// MatchUserName 匹配用户名
func MatchUserName(userName string) []string {
	matcher := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{4,15}$`)
	return matcher.FindAllString(userName, -1)
}
