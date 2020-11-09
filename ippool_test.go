package ippool

import (
	"bytes"
	"net"
	"testing"
)

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
	_, IPNet, _ := net.ParseCIDR("172.16.10.96/26")
	err := RegisterPrefix(&pool, IPNet)
	_, IPNet, _ = net.ParseCIDR("172.16.10.96/26")
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

func TestRequestIP4_1(t *testing.T) {
	// test first assignment
	expected_addr, _, _ := net.ParseCIDR("192.168.0.1/24")
	pool := InitPrefixPool()
	_, IPNet, _ := net.ParseCIDR("192.168.0.0/24")
	RegisterPrefix(&pool, IPNet)
	ipaddr, _ := RequestIP(&pool, IPNet)
	if bytes.Compare(ipaddr, expected_addr[12:16]) != 0 {
		t.Error()
	}
}

func TestRequestIP6(t *testing.T) {
	expected_addr, _, _ := net.ParseCIDR("FE80::1/72")
	pool := InitPrefixPool()
	_, IPNet, _ := net.ParseCIDR("FE80::0/72")
	RegisterPrefix(&pool, IPNet)
	ipaddr, _ := RequestIP(&pool, IPNet)
	if bytes.Compare(ipaddr, expected_addr) != 0 {
		t.Error()
	}
}

func TestReleaseIP4(t *testing.T) {

}

func TestReleaseIP6(t *testing.T) {

}
