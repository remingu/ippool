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
	"fmt"
	"math"
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

	//mask := ipNet.Mask
	if len(ipNet.IP) == 4 {
		fmt.Println(1)
	} else if len(ipNet.IP) == 16 {
		fmt.Println(1)
	}
	return net.IP{}
}

// ipv6

func LastFreeAddress(ipNet *net.IPNet) net.IP {
	var ui64 uint64
	ui64 = 0
	ui64 -= 1
	var ipv6 IPv6
	var last net.IP
	if len(ipNet.IP) == 4 {
		last = LastFreeAddress4(ipNet)
	} else if len(ipNet.IP) == 16 {
		prefix_length := GetPrefixLength(ipNet)
		for i := 0; i < 8; i++ {
			ipv6.H = append(ipv6.H, ipNet.IP[i])

		}
		for i := 8; i < 16; i++ {
			ipv6.L = append(ipv6.L, ipNet.IP[i])
		}
		//low := binary.BigEndian.Uint64(ipv6.L)
		//high := binary.BigEndian.Uint64(ipv6.H)
		fmt.Println(prefix_length)
		host_part := 128 - prefix_length
		var max_hosts uint64
		max_hosts = 0
		if host_part == 64 {
			max_hosts = ui64
		} else if host_part == 63 {
			max_hosts = 9223372036854775808
		} else if host_part < 63 {
			max_hosts = uint64(math.Exp2(float64(prefix_length)))
		} else {
			// rabbit hole for later
		}
		fmt.Println(max_hosts)
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

func GetMaxHosts(ipNet *net.IPNet) uint64 {
	prefix := GetPrefixLength(ipNet)
	hostrange := 128 - prefix
	max_hosts := (math.Exp2(float64(hostrange))) - 1
	return uint64(max_hosts)
}
