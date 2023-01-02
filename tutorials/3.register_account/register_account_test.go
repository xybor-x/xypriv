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

package register_account_test

import (
	"fmt"

	"github.com/xybor-x/xypriv"
)

// User implements Subject interface.
type User struct {
	id   string
	role string
}

// Relation returns the privilege of user over another subject. See the 2nd
// tutorial for further details.
func (u User) Relation(ctx any, subject xypriv.Subject) xypriv.Relation {
	switch u.role {
	case "admin":
		return "admin"
	case "mod":
		return "moderator"
	}

	switch t := subject.(type) {
	case User:
		if t.id == u.id {
			return "self"
		}
	}

	return "anyone"
}

// AccountTable implements StaticResource interface.
type AccountTable struct{}

// Context returns the context of AccountTable.
func (a AccountTable) Context() any {
	return nil
}

// Owner returns the owner of AccountTable.
func (a AccountTable) Owner() xypriv.Subject {
	return nil
}

// Permission returns the permission set of AccountTable.
func (a AccountTable) Permission(action ...string) xypriv.AccessLevel {
	// Return NotSupport if action is invalid.
	if len(action) != 2 {
		return xypriv.NotSupport
	}

	// If the action contains many strings, it is called parameterized action.
	// The access level for creating a record in AccountTable is different, it
	// depends on which role of the created account.
	// For the usage of parameterized action, see Example func.
	switch action[0] {
	case "create":
		switch action[1] {
		case "admin":
			// For the role of admin, the access level is HighSecret, only Admin
			// privilege can create this role.
			return xypriv.HighSecret
		case "mod":
			// For role of mod, the access level is HighConfidential, both Admin
			// and Moderator privilege can create this role.
			return xypriv.HighConfidential
		case "user":
			// For role of normal user, the access level is Public, it means
			// anyone can create this role.
			return xypriv.Public
		}
	}

	return xypriv.NotSupport
}

func Example() {
	var userAdmin = User{id: "Admin", role: "admin"}
	var userModerator = User{id: "Mod", role: "mod"}
	var user = User{id: "User"}
	var account = AccountTable{}

	if xypriv.Check(userAdmin).Perform("create", "admin").On(account) == nil {
		fmt.Println("admin can create admin account")
	}

	if xypriv.Check(userModerator).Perform("create", "mod").On(account) == nil {
		fmt.Println("mod can't create mod account")
	}

	if xypriv.Check(userModerator).Perform("create", "admin").On(account) != nil {
		fmt.Println("mod can't create admin account")
	}

	if xypriv.Check(user).Perform("create", "mod").On(account) != nil {
		fmt.Println("user can't create mod account")
	}

	if xypriv.Check(nil).Perform("create", "user").On(account) == nil {
		fmt.Println("anyone can create user account")
	}

	// Output:
	// admin can create admin account
	// mod can't create mod account
	// mod can't create admin account
	// user can't create mod account
	// anyone can create user account
}
