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
	_ "fmt"
	"net"
	"sync"
)

type Prefix struct {
	mutex           *sync.Mutex
	Prefix          net.IPNet
	Used            uint64
	Released		[]RCont
	max_hosts       uint64
}

type RCont struct {
	pos uint64
	blocksize uint64
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
	// wm -> tbdn
	var ret_ip net.IP

	ref_pool := *pool_ref
	network := GetNetLiteral(prefix)
	pool := ref_pool[network]
	if pool.Used <= pool.max_hosts {
		if pool.Used == 0 {
			ret_ip = pool.first_available
			pool.Used += 1
		}
	}




	return (ret_ip, nil)

}

func InitPrefixPool() map[string]Prefix {
	pool := make(map[string]Prefix)
	return pool
}

func GetNextAddress(index uint64) uint64 {

}