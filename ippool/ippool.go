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
	Used            []net.IP
	first_available net.IP
	first           net.IP
	last            net.IP
	max_hosts       uint64
}

func RegisterPrefix(pool_ref *map[string]Prefix, prefix *net.IPNet) {
	var new_prefix string
	new_prefix = GetNetLiteral(prefix)
	pool := *pool_ref
	pool[new_prefix] = Prefix{}
	InitPrefix(pool_ref, prefix, new_prefix)
}

func InitPrefix(pool_ref *map[string]Prefix, prefix *net.IPNet, prefix_string string) {
	var mutex = &sync.Mutex{}
	mutex.Lock()
	ref_pool := *pool_ref
	pool := ref_pool[prefix_string]
	first_addr := FirstFreeAddress(prefix)
	if len(prefix.IP) == 4 {
		last_addr := LastFreeAddress(prefix)
		max_hosts, _ := GetMaxHosts(prefix)
		pool.last = last_addr
		pool.first = first_addr
		pool.first_available = first_addr
		pool.max_hosts = max_hosts
		ref_pool[prefix_string] = pool
	} else if len(prefix.IP) == 16 {
		max_hosts, _ := GetMaxHosts(prefix)
		pool.max_hosts = max_hosts
		pool.first = first_addr
		pool.first_available = first_addr
		last_addr := LastFreeAddress(prefix)
		pool.last = last_addr
		ref_pool[prefix_string] = pool
	}
	mutex.Unlock()
}

func RequestIP(pool_ref *map[string]Prefix, prefix *net.IPNet) net.IP {
	// wm -> tbdn
	mutex := &sync.Mutex{}
	mutex.Lock()
	ref_pool := *pool_ref
	network := GetNetLiteral(prefix)
	pool := ref_pool[network]
	ret_ip := pool.first_available
	mutex.Unlock()
	return ret_ip

}

func InitPrefixPool() map[string]Prefix {
	pool := make(map[string]Prefix)
	return pool
}
