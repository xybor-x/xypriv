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

import (
	"fmt"
	"reflect"
	"strings"
)

// Privilege represents for relation of a Subject over another one.
type Privilege int

// AccessLevel represents for access level of a Resource.
type AccessLevel int

// Resource instances are objects that need to be protected from illegal
// actions.
type Resource interface {
	// Context returns the context of resource.
	Context() any

	// Owner returns the owner of resource
	Owner() Subject
}

// SimpleResource instances are Resources whose permission is based on only the
// action.
type SimpleResource interface {
	Resource

	// Permission defines the access level of resource. It is only based on the
	// action.
	Permission(action ...string) AccessLevel
}

// EnhancedResource instances are Resources whose permission is based on both
// the action and the subject who wants to perform the action.
type EnhancedResource interface {
	Resource

	// Permission defines the access level of resource. It is based on both the
	// action and the subject who wants to perform the action.
	Permission(s Subject, action ...string) AccessLevel
}

// Subject instances are entities that want to perform action on Resource.
type Subject interface {
	// Relation returns the privilege value of the current Subject over passed
	// Subject.
	Relation(s Subject, ctx any) Privilege
}

// Checker supports check if a subject can perform action on resource or not.
type Checker struct {
	subject Subject
	action  []string
}

// Check returns a Checker with the subject.
func Check(s Subject) *Checker {
	return &Checker{subject: s}
}

// Perform adds action to Checker and returns itself.
func (c *Checker) Perform(action ...string) *Checker {
	c.action = action
	return c
}

// On checks if a subject can perform action on resource or not.
func (c *Checker) On(r Resource) error {
	var accessLevel AccessLevel
	switch t := r.(type) {
	case SimpleResource:
		accessLevel = t.Permission(c.action...)
	case EnhancedResource:
		accessLevel = t.Permission(c.subject, c.action...)
	default:
		panic(NotImplementedError.New(
			"expected an object implementing Resource or EnhancedResource"))
	}

	var ctx = r.Context()
	var owner = r.Owner()

	if ctx == owner && owner != nil {
		panic(XyprivError.New("do not use the owner as the context, you " +
			"should set the context as nil in this case"))
	}

	var privilege = Anyone
	if c.subject != nil {
		privilege = c.subject.Relation(r.Owner(), ctx)
	}

	if int(privilege) < int(accessLevel) {
		return PermissionError.Newf(
			"%s do not have the permission to %s %s",
			getName(c.subject), strings.Join(c.action, "_"), getName(r))
	}

	return nil
}

// getName returns the name of object.
func getName(r any) string {
	if s, ok := r.(fmt.Stringer); ok {
		return s.String()
	}

	var name = reflect.TypeOf(r).Name()
	var result = make([]rune, 0, len(name))

	for _, c := range name {
		if c >= 'A' && c <= 'Z' {
			c = c - ('A' - 'a')
			if len(result) > 0 {
				result = append(result, '_')
			}
		}
		result = append(result, c)
	}

	return string(result)
}
