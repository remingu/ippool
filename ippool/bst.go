// Copyright (c) 2020 Puneeth S (puneeth8994 - https://github.com/puneeth8994)

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

import "errors"

type BstNode struct {
	val   uint64
	left  *BstNode
	right *BstNode
}

func (t *BstNode) Insert(value uint64) error {
	if t == nil {
		return errors.New("Tree is nil")
	}
	if t.val == value {
		return errors.New("This node value already exists")
	}
	if t.val > value {
		if t.left == nil {
			t.left = &BstNode{val: value}
			return nil
		}
		return t.left.Insert(value)
	}
	if t.val < value {
		if t.right == nil {
			t.right = &BstNode{val: value}
			return nil
		}
		return t.right.Insert(value)
	}
	return nil
}

func (t *BstNode) Find(value uint64) (BstNode, bool) {
	if t == nil {
		return BstNode{}, false
	}
	switch {
	case value == t.val:
		return *t, true
	case value < t.val:
		return t.left.Find(value)
	default:
		return t.right.Find(value)
	}
}

func (t *BstNode) Delete(value uint64) {
	t.remove(value)
}

func (t *BstNode) remove(value uint64) *BstNode {
	if t == nil {
		return nil
	}
	if value < t.val {
		t.left = t.left.remove(value)
		return t
	}
	if value > t.val {
		t.right = t.right.remove(value)
		return t
	}
	if t.left == nil && t.right == nil {
		t = nil
		return nil
	}
	if t.left == nil {
		t = t.right
		return t
	}
	if t.right == nil {
		t = t.left
		return t
	}
	smallestValOnRight := t.right
	for {
		if smallestValOnRight != nil && smallestValOnRight.left != nil {
			smallestValOnRight = smallestValOnRight.left
		} else {
			break
		}
	}
	t.val = smallestValOnRight.val
	t.right = t.right.remove(t.val)
	return t
}

func (t *BstNode) FindLast() uint64 {
	if t.right == nil {
		return t.val
	}
	return t.right.FindLast()
}
