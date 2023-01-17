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

func Example() {
	var userAdmin = User{id: "Admin", role: "admin"}
	var userModerator = User{id: "Mod", role: "mod"}
	var user = User{id: "User"}

	var accountTable = xypriv.AbstractResource("account_table")
	accountTable.SetPermission(xypriv.HighSecret, "create", "admin")
	accountTable.SetPermission(xypriv.HighConfidential, "create", "mod")
	accountTable.SetPermission(xypriv.Public, "create", "user")

	if xypriv.Check(userAdmin).Perform("create", "admin").On(accountTable) == nil {
		fmt.Println("admin can create admin account")
	}

	if xypriv.Check(userModerator).Perform("create", "mod").On(accountTable) == nil {
		fmt.Println("mod can't create mod account")
	}

	if xypriv.Check(userModerator).Perform("create", "admin").On(accountTable) != nil {
		fmt.Println("mod can't create admin account")
	}

	if xypriv.Check(user).Perform("create", "mod").On(accountTable) != nil {
		fmt.Println("user can't create mod account")
	}

	if xypriv.Check(nil).Perform("create", "user").On(accountTable) == nil {
		fmt.Println("anyone can create user account")
	}

	// Output:
	// admin can create admin account
	// mod can't create mod account
	// mod can't create admin account
	// user can't create mod account
	// anyone can create user account
}
