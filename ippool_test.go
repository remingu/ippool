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

func TestRegisterPrefix6_DuplicatePrefixes(t *testing.T) {
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
	// test assignment
	expected_addr, _, _ := net.ParseCIDR("192.168.0.1/24")
	pool := InitPrefixPool()
	_, IPNet, _ := net.ParseCIDR("192.168.0.0/24")
	RegisterPrefix(&pool, IPNet)
	ipaddr, _ := RequestIP(&pool, IPNet)
	if bytes.Compare(ipaddr, expected_addr[12:16]) != 0 {
		t.Error()
	}
}

func TestRequestIP4_2(t *testing.T) {
	// test if RequestIP() return error when prefix boundary is reached
	// last address is ipv4 broadcast and shall not be assigned
	pool := InitPrefixPool()
	_, IPNet, _ := net.ParseCIDR("192.168.0.0/28")
	RegisterPrefix(&pool, IPNet)
	for i := 0; i <= 14; i++ {
		addr, _ := RequestIP(&pool, IPNet)
		if i == 14 {
			if addr != nil {
				t.Error()
			}
		}
	}
}

func TestRequestIP6_1(t *testing.T) {
	//test assignment
	expected_addr, _, _ := net.ParseCIDR("FE80::1/72")
	pool := InitPrefixPool()
	_, IPNet, _ := net.ParseCIDR("FE80::0/72")
	RegisterPrefix(&pool, IPNet)
	ipaddr, _ := RequestIP(&pool, IPNet)
	if bytes.Compare(ipaddr, expected_addr) != 0 {
		t.Error()
	}
}

func TestRequestIP6_2(t *testing.T) {
	// test if RequestIP() returns error when prefix boundary is reached
	var i uint64
	pool := InitPrefixPool()
	_, IPNet, _ := net.ParseCIDR("FE80::0/120")
	RegisterPrefix(&pool, IPNet)
	for i = 0; i < Exp2nUInt64(8); i++ {
		ipaddr, _ := RequestIP(&pool, IPNet)
		if i == Exp2nUInt64(8)-1 && ipaddr != nil {
			t.Error()
		}
	}

}

func TestReleaseIP4(t *testing.T) {

}

func TestReleaseIP6(t *testing.T) {

}
