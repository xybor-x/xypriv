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

package context_test

import (
	"fmt"

	"github.com/xybor-x/xypriv"
)

// Group implements Subject interface
type Group struct {
	member []User
}

// Relation returns the privilege of group over another subject.
func (g Group) Relation(subject xypriv.Subject, ctx any) xypriv.Privilege {
	// A group doesn't have any privilege over other Subjects.
	return xypriv.Anyone
}

// IsMember checks if a user is a member of the group.
func (g Group) IsMember(u User) bool {
	for i := range g.member {
		if g.member[i].id == u.id {
			return true
		}
	}
	return false
}

// User implements Subject interface.
type User struct {
	id string
}

// Relation returns the privilege of user over another subject.
func (u User) Relation(subject xypriv.Subject, ctx any) xypriv.Privilege {
	// Self and other context-independent privilege should be evaluated first.
	switch t := subject.(type) {
	case User:
		if u.id == t.id {
			return xypriv.Self
		}
	}

	switch ctx.(type) {
	case nil: // The normal context.
		switch t := subject.(type) {
		case Group:
			if t.IsMember(u) {
				return xypriv.LowFamiliar
			}
		}
	case Group: // The group context.
		var g = ctx.(Group)

		// If the current user doesn't join in this group, the user has no
		// relation with group.
		if !g.IsMember(u) {
			return xypriv.Anyone
		}

		switch t := subject.(type) {
		case User:
			// The current user has no relation with subject user in this group,
			// if the subject doesn't join in group.
			if !g.IsMember(t) {
				return xypriv.Anyone
			}

			// In other cases, the current user and subject user joined in the
			// same group.
			return xypriv.LowFamiliar
		}
	}
	return xypriv.Anyone
}

// GroupPost implements Resource interface.
type GroupPost struct {
	user  User
	group Group
}

// Context returns the context of GroupPost.
func (gp GroupPost) Context() any {
	return gp.group
}

// Owner returns the owner of GroupPost.
func (gp GroupPost) Owner() xypriv.Subject {
	return gp.user
}

// Permission returns the permission set of GroupPost.
func (gp GroupPost) Permission(action ...string) xypriv.AccessLevel {
	// Return NotSupport if action is invalid.
	if len(action) != 1 {
		return xypriv.NotSupport
	}

	switch action[0] {
	case "update":
		// Admin and self privilege can update the post.
		return xypriv.HighSecret
	case "read":
		// Members in group can read the post.
		return xypriv.LowPrivate
	}
	return xypriv.NotSupport
}

func Example() {
	var userA = User{id: "A"}
	var userB = User{id: "B"}
	var userC = User{id: "C"}
	var group = Group{member: []User{userA, userB}}
	var post = GroupPost{user: userA, group: group}

	if xypriv.Check(userA).Perform("update").On(post) == nil {
		fmt.Println("userA can update the group post")
	}

	if xypriv.Check(userB).Perform("update").On(post) != nil {
		fmt.Println("userB can't update the group post")
	}

	if xypriv.Check(userC).Perform("update").On(post) != nil {
		fmt.Println("userC can't update the group post")
	}

	if xypriv.Check(userA).Perform("read").On(post) == nil {
		fmt.Println("userA can read the group post")
	}

	if xypriv.Check(userB).Perform("read").On(post) == nil {
		fmt.Println("userB can read the group post")
	}

	if xypriv.Check(userC).Perform("read").On(post) != nil {
		fmt.Println("userC can't read the group post")
	}

	// Output:
	// userA can update the group post
	// userB can't update the group post
	// userC can't update the group post
	// userA can read the group post
	// userB can read the group post
	// userC can't read the group post
}
