package tools

import "regexp"

// ReContent 获取部分文章内容
func ReContent(s string) string {
	reg := []map[string]string{
		{"reg": "<[^>]+>", "value": ""},
		{"reg": "&nbsp;", "value": ""},
	}
	content := s
	for _, item := range reg {
		re := regexp.MustCompile(item["reg"])
		content = re.ReplaceAllString(content, item["value"])
	}
	return Substr(content, 0, 500) + "..."
}

// Substr 字符串截取
func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}
