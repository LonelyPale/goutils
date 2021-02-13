package goutils

import "unicode"

// 是否首字母大写
func IsStartUpper(s string) bool {
	return unicode.IsUpper([]rune(s)[0])
}

// 是否首字母小写
func IsStartLower(s string) bool {
	return unicode.IsLower([]rune(s)[0])
}

// 首字母大写
func StartUpper(s string) string {
	if IsStartUpper(s) {
		return s
	}
	rs := []rune(s)
	rs[0] -= 32 //string的码表相差32位
	return string(rs)
}

// 首字母小写
func StartLower(s string) string {
	if IsStartLower(s) {
		return s
	}
	rs := []rune(s)
	rs[0] += 32 //string的码表相差32位
	return string(rs)
}
