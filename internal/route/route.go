package route

type Route struct {
	Protocol string
	Listen   string
	Upstream string
	AllowASN []int // ASN whitelist  empty = no filter
	DenyASN  []int // ASN blacklist  only used when AllowASN is empty
}
