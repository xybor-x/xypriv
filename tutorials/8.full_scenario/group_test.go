// MIT License
//
// Copyright (c) 2022 xybor-x
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package full_scenario_test

import (
	"github.com/xybor-x/xypriv"
)

// Group implements Subject interface
type Group struct {
	members    []*User
	admin      *User
	bannedUser []*User
}

// Relation returns the privilege of group over another subject.
func (g *Group) Relation(subject xypriv.Subject, ctx any) xypriv.Privilege {
	// A group should not have any privilege over other Subjects.
	return xypriv.Anyone
}

// IsMember checks if a user is a member of the group.
func (g Group) IsMember(u User) bool {
	for i := range g.members {
		if g.members[i].id == u.id {
			return true
		}
	}
	return false
}

// IsBanned checks if a user is banned by the group.
func (g Group) IsBanned(u User) bool {
	for i := range g.bannedUser {
		if g.bannedUser[i].id == u.id {
			return true
		}
	}
	return false
}
