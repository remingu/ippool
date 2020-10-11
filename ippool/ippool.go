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
	"encoding/binary"
	"errors"
	"net"
	"sync"
)

type Prefix struct {
	mutex      *sync.Mutex
	Prefix     net.IPNet
	Used       uint64
	Released32 []Cont32
	Released64 []Cont64
	max_hosts  uint64
}

type Cont32 struct {
	pos        uint32
	addr_range uint32
}

type Cont64 struct {
	pos        uint64
	addr_range uint64
}

func RegisterPrefix(pool_ref *map[string]Prefix, prefix *net.IPNet) {
	var new_prefix string
	new_prefix = GetNetLiteral(prefix)
	pool := *pool_ref
	pool[new_prefix] = Prefix{}
	InitPrefix(pool_ref, prefix, new_prefix)
}

func InitPrefix(pool_ref *map[string]Prefix, prefix *net.IPNet, prefix_string string) {
	ref_pool := *pool_ref
	pool := ref_pool[prefix_string]
	if len(prefix.IP) == 4 {
		max_hosts, _ := GetMaxHosts(prefix)
		pool.max_hosts = max_hosts
		pool.Used = 0
		ref_pool[prefix_string] = pool
	} else if len(prefix.IP) == 16 {
		max_hosts, _ := GetMaxHosts(prefix)
		pool.max_hosts = max_hosts
		pool.Used = 0
	}
	ref_pool[prefix_string] = pool

}

func RequestIP(pool_ref *map[string]Prefix, prefix *net.IPNet) (net.IP, error) {
	var IPaddr net.IP
	ref_pool := *pool_ref
	network := GetNetLiteral(prefix)
	pool := ref_pool[network]
	if pool.Used < pool.max_hosts {
		IPaddr = GetNextAddress(prefix, pool.Used)
		pool.Used += 1
		ref_pool[network] = pool
	} else {
		return nil, errors.New("prefix is full")
	}
	return IPaddr, nil
}

func InitPrefixPool() map[string]Prefix {
	pool := make(map[string]Prefix)
	return pool
}

func GetNextAddress(prefix *net.IPNet, index uint64) net.IP {
	var IPAddr net.IP
	if len(prefix.IP) == 4 {
		i := binary.BigEndian.Uint32(prefix.IP)
		i += 1 + uint32(index)
		new_addr := make([]byte, 4)
		binary.BigEndian.PutUint32(new_addr, i)
		IPAddr = new_addr
		return IPAddr
	} else {
		addr := GetIpv6Struct(prefix)
		i := binary.BigEndian.Uint64(addr.L)
		i += 1 + index
		new_addr := make([]byte, 8)
		binary.BigEndian.PutUint64(new_addr, i)
		IPAddr = append(IPAddr, addr.H...)
		IPAddr = append(IPAddr, new_addr...)
		return IPAddr
	}
}
