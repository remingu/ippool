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
	"math/bits"
	"net"
	"strconv"
)

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
	network := ipNet.IP
	var fa_ip net.IP
	fa_ip = network
	for i := len(fa_ip) - 1; i >= 0; i-- {
		if fa_ip[i]+1 <= 255 {
			fa_ip[i] = fa_ip[i] + 1
			break
		} else {
			fa_ip[i] = fa_ip[i] + 1
		}
	}
	return fa_ip
}

// ipv6

func LastFreeAddress(ipNet *net.IPNet) net.IP {
	var last net.IP
	if len(ipNet.IP) == 4 {
		last = LastFreeAddress4(ipNet)
	} else if len(ipNet.IP) == 16 {

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
