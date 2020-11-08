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

func Exp2nUInt64(a int) uint64 {
	// will give (2^64)-1 for 2^64 due to memory boundary - fine for us
	var i uint64
	var b uint64
	b = 1
	for i = 0; i < (uint64(a)); i++ {
		if i < 63 {
			b = b * 2
		} else if i == 63 {
			b = (b * 2) - 1
		}
	}
	return b
}
