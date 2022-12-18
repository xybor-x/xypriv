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

// GroupAvatar implements Resource interface.
type GroupAvatar struct {
	group *Group
}

// Context returns the context of GroupAvatar.
func (gp GroupAvatar) Context() any {
	// Group avatar itself belongs to Group, so it should not in Group context.
	return nil
}

// Owner returns the owner of GroupAvatar.
func (gp GroupAvatar) Owner() xypriv.Subject {
	return gp.group
}

// Permission returns the permission set of GroupAvatar.
func (gp GroupAvatar) Permission(action ...string) xypriv.AccessLevel {
	// Return NotSupport if action is invalid.
	if len(action) != 1 {
		return xypriv.NotSupport
	}

	switch action[0] {
	case "update":
		// Only admin or group admin privilege can update the post.
		return xypriv.LowSecret
	case "delete":
		// Only admin or group admin privilege can delete the post.
		return xypriv.LowSecret
	case "read":
		// Anyone can read the group avatar.
		return xypriv.Public
	}
	return xypriv.NotSupport
}
