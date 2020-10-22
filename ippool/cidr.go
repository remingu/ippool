// Copyright (c) 2020 Doc.ai and/or its affiliates.
// Copyright (c) 2020 Daniel Schlifka

// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package ipreg provides common functions useful when working with Classless Inter-Domain Routing (CIDR)

package ippool

import (
	"encoding/binary"
	"errors"
	"math/bits"
	"net"
	"strconv"
)

type IPv6 struct {
	H []byte
	L []byte
}

// ipv4

// NetworkAddress teturns the first IP address of an IP network
func NetworkAddress(ipNet *net.IPNet) net.IP {
	prefixNetwork := ipNet.IP.Mask(ipNet.Mask)
	return prefixNetwork
}

// BroadcastAddress returns the last IP address of an IP network
func BroadcastAddress4(ipNet *net.IPNet) net.IP {
	first := NetworkAddress(ipNet)
	ones, bit := ipNet.Mask.Size()
	var shift uint32 = 1
	shift <<= bit - ones
	intip := binary.BigEndian.Uint32(first.To4())
	intip = intip + shift - 1
	last := make(net.IP, 4)
	binary.BigEndian.PutUint32(last, intip)
	return last
}

func LastFreeAddress4(ipNet *net.IPNet) net.IP {
	first := NetworkAddress(ipNet)
	ones, bit := ipNet.Mask.Size()
	var shift uint32 = 1
	shift <<= bit - ones
	intip := binary.BigEndian.Uint32(first.To4())
	intip = intip + shift - 2
	last := make(net.IP, 4)
	binary.BigEndian.PutUint32(last, intip)
	return last
}

func FirstFreeAddress(ipNet *net.IPNet) net.IP {
	address := net.IP{}
	if len(ipNet.IP) == 4 {
		start_index := binary.BigEndian.Uint32(ipNet.IP)
		start_index++
		ip := make([]byte, 4)
		binary.BigEndian.PutUint32(ip, start_index)
		address = ip
	} else {
		ipv6 := GetIpv6Struct(ipNet)
		start_index_low := binary.BigEndian.Uint64(ipv6.L)
		ip := make([]byte, 8)
		start_index_low++
		binary.BigEndian.PutUint64(ip, start_index_low)
		ipv6.L = ip
		var addr []byte
		addr = append(addr, ipv6.H...)
		addr = append(addr, ip...)
		address = addr
	}
	return address
}

// ipv6

func LastFreeAddress(ipNet *net.IPNet) net.IP {
	var last net.IP
	var stop_index uint64
	if len(ipNet.IP) == 4 {
		last = LastFreeAddress4(ipNet)
	} else if len(ipNet.IP) == 16 {
		ipv6 := GetIpv6Struct(ipNet)
		max_hosts, _ := GetMaxHosts(ipNet)
		start_index := binary.BigEndian.Uint64(ipv6.L)
		if GetPrefixLength(ipNet) == 64 {
			stop_index = start_index + max_hosts + 1
		} else {
			stop_index = start_index + max_hosts
		}
		ip := make([]byte, 8)
		binary.BigEndian.PutUint64(ip, stop_index)
		ipv6.L = ip
		var addr []byte
		addr = append(addr, ipv6.H...)
		addr = append(addr, ip...)
		last = addr
	}

	return last
}

func GetPrefixLength(ipNet *net.IPNet) int {
	var prefix_length int
	for _, octet := range ipNet.Mask {
		if bits.TrailingZeros8(uint8(octet)) == 0 {
			prefix_length += 8
		} else if bits.TrailingZeros8(uint8(octet)) > 0 {
			prefix_length += 8 - bits.TrailingZeros8(uint8(octet))
		}
	}
	return prefix_length
}

func GetNetLiteral(prefix *net.IPNet) string {
	network := net.IP.String(NetworkAddress(prefix))
	prefix_length := GetPrefixLength(prefix)
	nw_string := network + "/" + strconv.Itoa(prefix_length)
	return nw_string
}

func GetMaxHosts(ipNet *net.IPNet) (uint64, error) {
	//check needed to exclude Multicast/Anycast address ranges
	var max_hosts uint64
	var prefix_length int
	if len(ipNet.IP) == 4 {
		prefix_length = GetPrefixLength(ipNet)
		if prefix_length >= 0 {
			max_hosts = Exp2nUInt64(32-prefix_length) - 2
			if max_hosts < 0 {
				return 0, errors.New("invalid prefix")
			}
		}
	} else {
		prefix_length = 128 - GetPrefixLength(ipNet)
		max_hosts = Exp2nUInt64(prefix_length) - 1
		if prefix_length > 64 {
			return 0, errors.New("invalid prefix size - must be between 64-128")
		}
	}
	return max_hosts, nil
}

func GetIpv6Struct(ipNet *net.IPNet) IPv6 {
	var ipv6 IPv6
	for i := 0; i < 8; i++ {
		ipv6.H = append(ipv6.H, ipNet.IP[i])

	}
	for i := 8; i < 16; i++ {
		ipv6.L = append(ipv6.L, ipNet.IP[i])
	}
	return ipv6
}
