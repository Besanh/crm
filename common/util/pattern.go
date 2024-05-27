package util

import (
	"regexp"
	"strings"
)

func CheckPatternTelephone(telelphone string) bool {
	// pattern := "^[1]\\d{9}|[9]\\d{8}|[2-8]\\d{8,9}$"
	r, _ := regexp.Compile(`^(0?|\\+84|\\+840)\\d{9}$`)
	return r.MatchString(telelphone)
}

func CheckPatternOTP(otpVoice string) bool {
	// pattern := "^[1]\\d{9}|[9]\\d{8}|[2-8]\\d{8,9}$"
	r, _ := regexp.Compile(`^\\d+$`)
	return r.MatchString(otpVoice)
}

func CheckPatternNumeric(Key string) bool {
	// pattern := "^[1]\\d{9}|[9]\\d{8}|[2-8]\\d{8,9}$"
	r, _ := regexp.Compile(`^\\d+$`)
	return r.MatchString(Key)
}
func CheckPatternNumericText(Key string) bool {
	r, _ := regexp.Compile(`^#\\d+$`)
	return r.MatchString(Key)
}

func Replace84To0(str string) string {
	r := regexp.MustCompile(`^(\\+84|84)`)
	return r.ReplaceAllString(str, "0")
}

func TrimPrefixPhone(telelphone string) string {
	telelphone = strings.TrimPrefix(telelphone, "0")
	telelphone = strings.TrimPrefix(telelphone, "840")
	telelphone = strings.TrimPrefix(telelphone, "+840")
	telelphone = strings.TrimPrefix(telelphone, "+84")
	telelphone = strings.TrimPrefix(telelphone, "84")
	return telelphone
}
