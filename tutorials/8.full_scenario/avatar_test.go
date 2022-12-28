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

// Permission returns the permission set of Avatar.
func (a Avatar) Permission(action ...string) xypriv.AccessLevel {
	if len(action) != 1 {
		return xypriv.NotSupport
	}

	switch action[0] {
	case "update":
		return xypriv.HighSecret
	case "comment":
		return xypriv.LowPrivate
	case "read":
		return xypriv.Public
	}
	return xypriv.NotSupport
}
