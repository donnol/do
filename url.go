package do

import (
	"net/url"
	"strings"
)

// ReplaceIP replace link 's ip with nip
func ReplaceIP(link, ip, nip string) (r string, err error) {
	info, err := url.Parse(link)
	if err != nil {
		return
	}
	if info.Host == "" {
		r = strings.ReplaceAll(link, ip, nip)
		return
	}
	host := strings.ReplaceAll(info.Host, ip, nip)
	info.Host = host
	r = info.String()
	return
}
