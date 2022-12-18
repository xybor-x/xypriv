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

package banned_user_test

import (
	"fmt"

	"github.com/xybor-x/xypriv"
)

// User implements Subject interface.
type User struct {
	id         string
	bannedUser []*User
}

// Relation returns the privilege of user over another subject.
func (u *User) Relation(subject xypriv.Subject, ctx any) xypriv.Privilege {
	switch t := subject.(type) {
	case *User:
		// Check if subject is the current user or not.
		if u.id == t.id {
			return xypriv.Self
		}
	}

	switch ctx {
	case nil:
		switch t := subject.(type) {
		case *User:
			// Check if the current user is in banned list of subject.
			for i := range t.bannedUser {
				if t.bannedUser[i].id == u.id {
					return xypriv.BadRelation
				}
			}
		}
	}
	return xypriv.Anyone
}

// Avatar implements Resource interface.
type Avatar struct {
	user *User
}

// Context returns the context of Avatar.
func (a Avatar) Context() any {
	return nil
}

// Owner returns the owner of Avatar.
func (a Avatar) Owner() xypriv.Subject {
	return a.user
}

// Permission returns the permission set of Avatar. See the first tutorial for
// further details.
func (a Avatar) Permission(action ...string) xypriv.AccessLevel {
	if len(action) != 1 {
		return xypriv.NotSupport
	}

	switch action[0] {
	case "read":
		return xypriv.Public
	}

	return xypriv.NotSupport
}

func Example() {
	var userA = &User{id: "A"}
	var userB = &User{id: "B"}
	var userC = &User{id: "C"}

	userA.bannedUser = append(userA.bannedUser, userB)
	var avtA = Avatar{user: userA}

	if xypriv.Check(userA).Perform("read").On(avtA) == nil {
		fmt.Println("userA can read avtA")
	}

	if xypriv.Check(userB).Perform("read").On(avtA) != nil {
		fmt.Println("userB can't read avtA")
	}

	if xypriv.Check(userC).Perform("read").On(avtA) == nil {
		fmt.Printf("userC can read avtA")
	}

	// Output:
	// userA can read avtA
	// userB can't read avtA
	// userC can read avtA
}
