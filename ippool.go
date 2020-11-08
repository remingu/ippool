// Copyright (c) 2020 Daniel Schlifka
//
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
package ippool

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
)

type Prefix struct {
	Prefix      net.IPNet
	Used        uint64
	FreedIPs    uint64
	ReleasedIPs BstNode
	max_hosts   uint64
}

func RegisterPrefix(pool_ref *map[string]Prefix, prefix *net.IPNet) {
	// registers a new prefix to prefix map.
	var new_prefix string
	new_prefix = GetNetLiteral(prefix)
	pool := *pool_ref
	pool[new_prefix] = Prefix{}
	InitPrefix(pool_ref, prefix, new_prefix)
}

func InitPrefix(pool_ref *map[string]Prefix, prefix *net.IPNet, prefix_string string) {
	// detect and set basic parameters for a prefix
	ref_pool := *pool_ref
	pool := ref_pool[prefix_string]
	max_hosts, _ := GetMaxHosts(prefix)
	pool.max_hosts = max_hosts
	pool.Used = 0
	pool.FreedIPs = 0
	ref_pool[prefix_string] = pool
}

func RequestIP(pool_ref *map[string]Prefix, prefix *net.IPNet) (net.IP, error) {
	/*

	 */
	var IPAddr net.IP
	ref_pool := *pool_ref
	network := GetNetLiteral(prefix)
	pool := ref_pool[network]
	if pool.Used < pool.max_hosts {
		if pool.FreedIPs == 0 {
			IPAddr = GetNextAddress(prefix, pool.Used)
			pool.Used++
			ref_pool[network] = pool
		} else {
			if len(prefix.IP) == 4 {
				addr_uint32 := uint32(pool.ReleasedIPs.FindLast())
				addr := make([]byte, 4)
				binary.BigEndian.PutUint32(addr, addr_uint32)
				pool.ReleasedIPs.Delete(uint64(addr_uint32))
				pool.FreedIPs--
				IPAddr = addr
			} else {
				addr := GetIpv6Struct(prefix)
				addr_uint64 := pool.ReleasedIPs.FindLast()
				new_addr := make([]byte, 8)
				binary.BigEndian.PutUint64(new_addr, addr_uint64)
				IPAddr = append(IPAddr, addr.H...)
				IPAddr = append(IPAddr, new_addr...)
				pool.ReleasedIPs.Delete(addr_uint64)
				pool.FreedIPs--
			}
		}
	} else {
		if pool.FreedIPs == 0 {
			return nil, errors.New("prefix is full")
		} else {
			if len(prefix.IP) == 4 {
				addr_uint32 := uint32(pool.ReleasedIPs.FindLast())
				addr := make([]byte, 4)
				binary.BigEndian.PutUint32(addr, addr_uint32)
				pool.ReleasedIPs.Delete(uint64(addr_uint32))
				pool.FreedIPs--
				IPAddr = addr
			} else {
				addr := GetIpv6Struct(prefix)
				addr_uint64 := pool.ReleasedIPs.FindLast()
				new_addr := make([]byte, 8)
				binary.BigEndian.PutUint64(new_addr, addr_uint64)
				IPAddr = append(IPAddr, addr.H...)
				IPAddr = append(IPAddr, new_addr...)
				pool.ReleasedIPs.Delete(addr_uint64)
				pool.FreedIPs--
			}
		}
	}
	ref_pool[network] = pool
	return IPAddr, nil
}

func InitPrefixPool() map[string]Prefix {
	// initializes a map which will contain all prefixes later on
	pool := make(map[string]Prefix)
	return pool
}

func GetNextAddress(prefix *net.IPNet, index uint64) net.IP {
	var IPAddr net.IP
	if len(prefix.IP) == 4 {
		i := binary.BigEndian.Uint32(prefix.IP)
		i += uint32(index) + 1
		new_addr := make([]byte, 4)
		binary.BigEndian.PutUint32(new_addr, i)
		IPAddr = new_addr
		return IPAddr
	} else {
		addr := GetIpv6Struct(prefix)
		i := binary.BigEndian.Uint64(addr.L)
		i += index + 1
		new_addr := make([]byte, 8)
		binary.BigEndian.PutUint64(new_addr, i)
		IPAddr = append(IPAddr, addr.H...)
		IPAddr = append(IPAddr, new_addr...)
		return IPAddr
	}
}

func ReleaseIP(pool_ref *map[string]Prefix, prefix *net.IPNet, addr net.IP) error {

	ref_pool := *pool_ref
	network := GetNetLiteral(prefix)
	pool := ref_pool[network]
	if len(prefix.IP) == 4 {
		addr_uint32 := binary.BigEndian.Uint32(addr.To4())
		addr_first := binary.BigEndian.Uint32(FirstFreeAddress(prefix))
		addr_last := binary.BigEndian.Uint32(LastFreeAddress(prefix))
		addr_last_used := (uint32(pool.Used) + addr_first)

		if addr_uint32 > addr_last || addr_uint32 < addr_first {
			// detect the freed IP really belongs to the registered prefix range
			return errors.New("invalid ip - address not in prefix range")
		} else if addr_uint32 > addr_last_used {
			/* we need to check that an IP is already assigned before adding it to the "free IPs" binary tree.
			   Otherwise it can end in duplicate address assignments.
			*/
			return errors.New("ip is not in use and must not be released")
		} else {
			// check
			pool.FreedIPs++
			x := make([]byte, 4)
			binary.BigEndian.PutUint32(x, addr_uint32)
			err := pool.ReleasedIPs.Insert(uint64(addr_uint32))
			if err != nil {
				return errors.New("unable to insert ip4")
			}
		}
	} else {
		prefixLength := GetPrefixLength(prefix)
		var AddrH, AddrL []byte
		for i := 0; i < 8; i++ {
			AddrH = append(AddrH, addr[i])
		}
		for i := 8; i < 16; i++ {
			AddrL = append(AddrL, addr[i])
		}
		network := GetIpv6Struct(prefix)
		// we support only /64 therefore upper 8bytes should always remain the same.
		if bytes.Compare(network.H, AddrH) != 0 {
			return errors.New("invalid ip - address not in prefix range")
		}
		if prefixLength == 64 {
			pool.FreedIPs++
			err := pool.ReleasedIPs.Insert(binary.BigEndian.Uint64(AddrL))
			if err != nil {
				return errors.New("unable to insert ip6")
			}
		} else {
			// check the ip is bigger than upper prefix boundary
			if binary.BigEndian.Uint64(AddrL) > binary.BigEndian.Uint64(network.L)+pool.max_hosts {
				return errors.New("invalid ip - address not in prefix range")
			}
			// check the ip is smaller than lower prefix boundary
			if binary.BigEndian.Uint64(AddrL) < binary.BigEndian.Uint64(network.L) {
				return errors.New("invalid ip - address not in prefix range")
			}
			/* we need to check that an IP is already assigned before adding it to the "free IPs" binary tree.
			   Otherwise it can end in duplicate address assignments. tbd
			*/
			if binary.BigEndian.Uint64(AddrL) > binary.BigEndian.Uint64(network.L)+pool.Used {
				return errors.New("ip is not in use and must not be released")
			}
			pool.FreedIPs++
			err := pool.ReleasedIPs.Insert(binary.BigEndian.Uint64(AddrL))
			if err != nil {
				return errors.New("unable to insert ip6")
			}
		}
	}
	ref_pool[network] = pool
	return nil
}

func IsIPInUse(pool_ref *map[string]Prefix, prefix *net.IPNet, addr net.IP) bool {
	ref_pool := *pool_ref
	network := GetNetLiteral(prefix)
	pool := ref_pool[network]
	i := binary.BigEndian.Uint32(addr.To4())
	_, val := pool.ReleasedIPs.Find(uint64(i))
	return val
}
