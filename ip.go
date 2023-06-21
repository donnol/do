package do

import "net/netip"

// IsValidIP check the ip if is a valid ipv4 or ipv6 addr.
func IsValidIP(ip string) bool {
	_, err := netip.ParseAddr(ip)
	return err == nil
}
