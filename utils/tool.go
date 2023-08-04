package utils

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	IPV4       = `^(((\d{1,2})|(1\d{2})|(2[0-4]\d)|(25[0-5]))\.){3}((\d{1,2})|(1\d{2})|(2[0-4]\d)|(25[0-5]))$`
	STATIC_KEY = `^[A-Za-z0-9!@#$%?&]{16}$`
)

func CheckIpv4(addstr string) bool {
	ipv4Regex, _ := regexp.Compile(IPV4)
	split := strings.Split(addstr, ":")
	if len(split) != 2 {
		return false
	}
	ip := split[0]
	port, err := strconv.Atoi(split[1])
	if err != nil {
		return false
	}
	if port <= 0 || port > 65535 {
		return false
	}
	return ipv4Regex.MatchString(ip)
}
func CheckKey(key string) bool {
	r := regexp.MustCompile(STATIC_KEY)
	return r.MatchString(key)
}
