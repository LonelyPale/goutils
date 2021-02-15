package goutils

import "unicode"

// 是否首字母大写
func IsStartUpper(s string) bool {
	if len(s) == 0 {
		return false
	}
	return unicode.IsUpper([]rune(s)[0])
}

// 是否首字母小写
func IsStartLower(s string) bool {
	if len(s) == 0 {
		return false
	}
	return unicode.IsLower([]rune(s)[0])
}

// 首字母大写
func StartUpper(s string) string {
	if len(s) == 0 {
		return ""
	}
	if IsStartUpper(s) {
		return s
	}
	rs := []rune(s)
	rs[0] -= 32 //string的码表相差32位
	return string(rs)
}

// 首字母小写
func StartLower(s string) string {
	if len(s) == 0 {
		return ""
	}
	if IsStartLower(s) {
		return s
	}
	rs := []rune(s)
	rs[0] += 32 //string的码表相差32位
	return string(rs)
}
