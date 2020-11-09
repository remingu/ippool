package ippool

import (
	"net"
	"testing"
)

func TestInitPrefixPool(t *testing.T) {

}

func TestRegisterPrefix4(t *testing.T) {
	// register ip4 prefix
	pool := InitPrefixPool()
	_, IPNet, _ := net.ParseCIDR("192.168.10.96/24")
	err := RegisterPrefix(&pool, IPNet)
	if err != nil {
		t.Error()
	}
}

func TestRegisterPrefix4_DuplicatePrefixes(t *testing.T) {
	// check for duplicate ipv4 prefixes
	pool := InitPrefixPool()
	_, IPNet, _ := net.ParseCIDR("192.168.10.96/24")
	err := RegisterPrefix(&pool, IPNet)
	_, IPNet, _ = net.ParseCIDR("192.168.10.96/24")
	err = RegisterPrefix(&pool, IPNet)
	if err == nil {
		t.Error()
	}
}

func TestRegisterPrefix6(t *testing.T) {
	// register ip6 prefix
	pool := InitPrefixPool()
	_, IPNet, _ := net.ParseCIDR("FE80::0/96")
	err := RegisterPrefix(&pool, IPNet)
	if err != nil {
		t.Error()
	}
}

func TestRegisterPrefix_DuplicatePrefixes(t *testing.T) {
	// check for duplicate ipv6 prefixes
	pool := InitPrefixPool()
	_, IPNet, _ := net.ParseCIDR("FE80::0/96")
	err := RegisterPrefix(&pool, IPNet)
	_, IPNet, _ = net.ParseCIDR("FE80::0/96")
	err = RegisterPrefix(&pool, IPNet)
	if err == nil {
		t.Error()
	}
}

func TestRequestIP4(t *testing.T) {

}

func TestRequestIP6(t *testing.T) {

}

func TestReleaseIP4(t *testing.T) {

}

func TestReleaseIP6(t *testing.T) {

}
