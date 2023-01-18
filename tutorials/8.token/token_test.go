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

package token_test

import (
	"fmt"

	"github.com/xybor-x/xypriv"
)

func init() {
	xypriv.AddRelation(nil, "banned", xypriv.BadRelation)
	xypriv.AddRelation(nil, "groupMember", xypriv.LowFamiliar)
	xypriv.AddRelation(Group{}, "anyone", xypriv.Anyone)
	xypriv.AddRelation(Group{}, "sameGroup", xypriv.LowFamiliar)
	xypriv.AddRelation(Group{}, "self", xypriv.Self)
}

// Group implements Subject interface
type Group struct {
	member []User
}

// Relation returns the privilege of group over another subject.
func (g Group) Relation(ctx any, subject xypriv.Subject) xypriv.Relation {
	// A group doesn't have any privilege over other Subjects.
	return "anyone"
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
func (u User) Relation(ctx any, subject xypriv.Subject) xypriv.Relation {
	// Self and other context-independent privilege should be evaluated first.
	switch t := subject.(type) {
	case User:
		if u.id == t.id {
			return "self"
		}
	}

	switch ctx.(type) {
	case nil: // The normal context.
		switch t := subject.(type) {
		case Group:
			if t.IsMember(u) {
				return "groupMember"
			}
		}
	case Group: // The group context.
		var g = ctx.(Group)

		// If the current user doesn't join in this group, the user has no
		// relation with group.
		if !g.IsMember(u) {
			return "anyone"
		}

		switch t := subject.(type) {
		case User:
			// The current user has no relation with subject user in this group,
			// if the subject doesn't join in group.
			if !g.IsMember(t) {
				return "anyone"
			}

			// In other cases, the current user and subject user joined in the
			// same group.
			return "sameGroup"
		}
	}
	return "anyone"
}

// GroupPost implements StaticResource interface.
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
	var user = User{id: "user"}

	var tokenReadPost = xypriv.NewToken()
	tokenReadPost.AllowActionInScope(GroupPost{}, "read")

	var tokenGroup = xypriv.NewToken()
	tokenGroup.AllowScope(Group{})

	var tokenNormalContext = xypriv.NewToken()
	tokenNormalContext.AllowScope(nil)

	var group = Group{member: []User{user}}
	var post = GroupPost{user: user, group: group}

	if xypriv.Check(user).Perform("update").On(post) == nil {
		fmt.Println("user can update the group post")
	}

	if xypriv.Check(user).Delegate(tokenReadPost).Perform("update").On(post) != nil {
		fmt.Println("tokenReadPost can't update the group post")
	}

	if xypriv.Check(user).Delegate(tokenGroup).Perform("update").On(post) == nil {
		fmt.Println("tokenGroup can update the group post")
	}

	if xypriv.Check(user).Delegate(tokenNormalContext).Perform("update").On(post) != nil {
		fmt.Println("tokenNormalContext can't update the group post")
	}

	if xypriv.Check(user).Perform("read").On(post) == nil {
		fmt.Println("user can read the group post")
	}

	if xypriv.Check(user).Delegate(tokenReadPost).Perform("read").On(post) == nil {
		fmt.Println("tokenReadPost can read the group post")
	}

	if xypriv.Check(user).Delegate(tokenGroup).Perform("read").On(post) == nil {
		fmt.Println("tokenGroup can read the group post")
	}

	if xypriv.Check(user).Delegate(tokenNormalContext).Perform("read").On(post) != nil {
		fmt.Println("tokenNormalContext can't read the group post")
	}

	// Output:
	// user can update the group post
	// tokenReadPost can't update the group post
	// tokenGroup can update the group post
	// tokenNormalContext can't update the group post
	// user can read the group post
	// tokenReadPost can read the group post
	// tokenGroup can read the group post
	// tokenNormalContext can't read the group post
}
