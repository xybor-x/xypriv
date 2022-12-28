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

package enhanced_resource_test

import (
	"fmt"

	"github.com/xybor-x/xypriv"
)

// User implements Subject interface.
type User struct {
	id string
}

// Relation returns the privilege of user over another subject.
func (u User) Relation(subject xypriv.Subject, ctx any) xypriv.Privilege {
	if ctx != nil {
		return xypriv.Anyone
	}

	// Check if subject is the current user or not.
	if u.id == subject.(User).id {
		return xypriv.Self
	}
	return xypriv.Anyone
}

// Avatar implements Resource interface.
type Avatar struct {
	user       User
	bannedUser []User
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
func (a Avatar) Permission(s xypriv.Subject, action ...string) xypriv.AccessLevel {
	// Return NotSupport if action is invalid.
	if len(action) != 1 {
		return xypriv.NotSupport
	}

	// Only support read as valid action. The access level for updating avatar
	// is Public. This means anyone can read the avatar.
	// But the avatar banned some users. If the subject is banned from the
	// resource, an access level of NotSupport will be returned.
	switch action[0] {
	case "read":
		if u, ok := s.(User); ok {
			for i := range a.bannedUser {
				if a.bannedUser[i].id == u.id {
					return xypriv.NotSupport
				}
			}
		}
		return xypriv.Public
	}
	return xypriv.NotSupport
}

func Example() {
	var userA = User{id: "A"}
	var userB = User{id: "B"}
	var userC = User{id: "C"}

	var avtA = Avatar{user: userA}
	avtA.bannedUser = append(avtA.bannedUser, userC)

	if xypriv.Check(userA).Perform("read").On(avtA) == nil {
		fmt.Println("userA can read avtA")
	}

	if xypriv.Check(userB).Perform("read").On(avtA) == nil {
		fmt.Println("userB can read avtA")
	}

	if xypriv.Check(userC).Perform("read").On(avtA) != nil {
		fmt.Printf("userC can't read avtA")
	}

	// Output:
	// userA can read avtA
	// userB can read avtA
	// userC can't read avtA
}
