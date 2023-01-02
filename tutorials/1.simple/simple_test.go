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

package simple_test

import (
	"fmt"

	"github.com/xybor-x/xypriv"
)

// User implements Subject interface.
type User struct {
	id string
}

// Relation returns the privilege of user over another subject.
func (u User) Relation(ctx any, subject xypriv.Subject) xypriv.Relation {
	// Check if subject is the current user or not.
	switch t := subject.(type) {
	case User:
		if u.id == t.id {
			return "self"
		}
	}
	return "anyone"
}

// Avatar implements StaticResource interface.
type Avatar struct {
	user User
}

// Context returns the context of Avatar.
func (a Avatar) Context() any {
	return nil
}

// Owner returns the owner of Avatar.
func (a Avatar) Owner() xypriv.Subject {
	return a.user
}

// Permission returns the permission set of Avatar.
func (a Avatar) Permission(action ...string) xypriv.AccessLevel {
	// Return NotSupport if action is invalid.
	if len(action) != 1 {
		return xypriv.NotSupport
	}

	// Only support update as valid action. The access level for updating avatar
	// is LowSecret. This means only privilege of admins or self can perform
	// this action.
	switch action[0] {
	case "update":
		return xypriv.TopSecret
	}
	return xypriv.NotSupport
}

func Example() {
	var userA = User{id: "A"}
	var userB = User{id: "B"}

	var avtA = Avatar{user: userA}

	if xypriv.Check(userA).Perform("update").On(avtA) == nil {
		fmt.Println("userA can update avtA")
	}

	if xypriv.Check(userB).Perform("update").On(avtA) != nil {
		fmt.Println("userB can't update avtA")
	}

	if xypriv.Check(userA).Perform("read").On(avtA) != nil {
		fmt.Printf("userA can't read avtA (because this action is not defined)")
	}

	// Output:
	// userA can update avtA
	// userB can't update avtA
	// userA can't read avtA (because this action is not defined)
}
