package utils

import (
	"strings"
	"unicode"
)

var vi = map[string]string{
	"gte":              "lớn hơn hoặc bằng",
	"gt":               "lớn hơn",
	"lte":              "nhỏ hơn hoặc bằng",
	"lt":               "nhỏ hơn",
	"in":               "nằm trong",
	"nin":              "không nằm trong",
	"passwordError":    "phải có ít nhất 8 ký tự, bao gồm ít nhất 1 chữ in hoa, 1 chữ số và 1 ký tự đặc biệt",
	"mustBe":           "phải",
	"characters":       "kí tự",
	"required":         "bắt buộc nhập",
	"record not found": "không tìm thấy dữ liệu",
	"user's status is": "trạng thái của người dùng đang",
	"active":           "hoạt động",
	"inactive":         "không hoạt động",
	"suspended":        "bị cấm",
	"token has invalid claims: token is expired":       "token không hợp lệ",
	"token signature is invalid: signature is invalid": "token không hợp lệ",
	"domain":        "tên miền",
	"availableList": "danh sách khả dụng",
	"Access denied": "từ chối truy cập",
}

var en = map[string]string{
	"gte":              "greater than or equal",
	"gt":               "greater than",
	"lte":              "less than or equal",
	"lt":               "less than",
	"in":               "in",
	"nin":              "not in",
	"passwordError":    "must be at least 8 characters including 1 upper case, 1 digit and 1 special character",
	"mustBe":           "must be",
	"characters":       "characters",
	"required":         "required",
	"record not found": "cannot find data",
	"user's status is": "user's status is",
	"active":           "active",
	"inactive":         "inactive",
	"suspended":        "suspended",
	"token has invalid claims: token is expired":       "token is invalid",
	"token signature is invalid: signature is invalid": "token is invalid",
	"domain":        "domain",
	"availableList": "available list",
	"Access denied": "access denied",
}

var Translation = map[string]map[string]string{
	"en_US": en,
	"vi_VN": vi,
}

func translate(k string, locale string) string {
	trm := map[string]string{}
	if len(locale) > 0 && Translation[locale] != nil {
		trm = Translation[locale]
	}

	if v, ok := trm[k]; ok {
		return v
	}
	return k
}

func Reason(locale string, reasons ...string) string {
	reason := ""
	for _, eachReason := range reasons {
		if eachReason == "" {
			break
		}
		reason += translate(eachReason, locale) + " "
	}

	reason = strings.TrimSpace(reason)
	if reason == "" {
		return reason
	}

	firstChar := unicode.ToUpper(rune(reason[0]))
	return string(firstChar) + reason[1:]
}
