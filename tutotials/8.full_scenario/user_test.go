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

// User implements Subject interface.
type User struct {
	id         string
	role       string
	bannedUser []*User
	friends    []*User
}

// IsFriend check if another user is a friend of the current user.
func (u User) IsFriend(another User) bool {
	for i := range u.friends {
		if u.friends[i].id == another.id {
			return true
		}
	}
	return false
}

// IsBanned check if another user is banned by the current user.
func (u User) IsBanned(another User) bool {
	for i := range u.bannedUser {
		if u.bannedUser[i].id == another.id {
			return true
		}
	}
	return false
}

// Relation returns the privilege of user over another subject.
func (u *User) Relation(subject xypriv.Subject, ctx any) xypriv.Privilege {
	// Non-context privilege will be served first.
	if usr, ok := subject.(*User); ok {
		if usr.id == u.id {
			return xypriv.Self
		}
	}

	// Role-based is also a non-context privilege.
	switch u.role {
	case "admin":
		return xypriv.Admin
	case "mod":
		return xypriv.Moderator
	}

	switch ctx.(type) {
	// In the normal context, a user has the following relations:
	// 1. u is in friend list of t (User).
	// 2. u is banned by t (User).
	// 3. u is admin of t (Group).
	// 4. u is member of t (Group).
	// 5. u is banned by t (Group)
	case nil:
		switch t := subject.(type) {
		case *User:
			if t.IsFriend(*u) {
				return xypriv.LowFamiliar
			}
			if t.IsBanned(*u) {
				return xypriv.BadRelation
			}
		case *Group:
			if u.id == t.admin.id {
				return xypriv.LocalAdmin
			}

			if t.IsMember(*u) {
				return xypriv.LowFamiliar
			}

			if t.IsBanned(*u) {
				return xypriv.BadRelation
			}
		}
	// In the group context:
	// 1. If u is not a member of group, or banned by the group, the privilege
	//    is Anyone or BadRelation, respectively.
	// 2. If subject is a user, and not a member of group, or banned by the
	//    group, the privilege is the same as above. Otherwise, the current user
	//    has the following relations:
	//    a. u is group admin of t.
	//    b. u is same group with t.
	//    c. u is banned by t.
	case *Group:
		var group = ctx.(*Group)
		if group.IsBanned(*u) {
			return xypriv.BadRelation
		}
		if !group.IsMember(*u) {
			return xypriv.Anyone
		}

		switch t := subject.(type) {
		case *User:
			if group.IsBanned(*t) {
				return xypriv.BadRelation
			}
			if !group.IsMember(*t) {
				return xypriv.Anyone
			}
			if group.admin.id == u.id {
				return xypriv.LocalAdmin
			}
			if t.IsBanned(*u) {
				return xypriv.BadRelation
			}
			return xypriv.LowFamiliar
		}
	}

	return xypriv.Anyone
}
