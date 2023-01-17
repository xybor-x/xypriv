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

package xypriv

import "strings"

//abstractResourceStore stores all created abstract resources.
var abstractResourceStore = map[string]AbstractResourceDetails{}

// AbstractResource returns an existed AbstractResourceDetails or creates one
// if it doesn't exist before.
func AbstractResource(name string) AbstractResourceDetails {
	if _, ok := abstractResourceStore[name]; !ok {
		abstractResourceStore[name] = AbstractResourceDetails{
			permissions: make(map[string]AccessLevel),
		}
	}

	return abstractResourceStore[name]
}

// AbstractResourceDetails contains information about an abstract resource.
type AbstractResourceDetails struct {
	permissions map[string]AccessLevel
	context     any
	owner       Subject
}

// SetContext sets the context of resource.
func (r *AbstractResourceDetails) SetContext(c any) {
	r.context = c
}

// SetOwner sets the owner of resource.
func (r *AbstractResourceDetails) SetOwner(o Subject) {
	r.owner = o
}

// SetPermission sets the access level corresponding to the action.
func (r *AbstractResourceDetails) SetPermission(l AccessLevel, action ...string) {
	r.permissions[strings.Join(action, "_")] = l
}

// Context implements Resource interface.
func (r AbstractResourceDetails) Context() any {
	return r.context
}

// Owner implements Resource interface.
func (r AbstractResourceDetails) Owner() Subject {
	return r.owner
}

// Permission implements Resource interface.
func (r AbstractResourceDetails) Permission(action ...string) AccessLevel {
	var actions = strings.Join(action, "_")
	if val, ok := r.permissions[actions]; ok {
		return val
	}
	return NotSupport
}
