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

// Relation represents for the relation of a Subject over another one.
type Relation string

// Privilege represents for the privilege of a Subject over another one.
type Privilege int

// AccessLevel represents for the access level of a Resource.
type AccessLevel int

// Resource instances are objects that need to be protected from illegal
// actions.
type Resource interface {
	// Context returns the context of resource.
	Context() any

	// Owner returns the owner of resource
	Owner() Subject
}

// StaticResource instances are Resources whose permission is based on only the
// action.
type StaticResource interface {
	Resource

	// Permission defines the access level of resource. It is only based on the
	// action.
	Permission(action ...string) AccessLevel
}

// DynamicResource instances are Resources whose permission is based on both
// the action and the subject who wants to perform the action.
type DynamicResource interface {
	Resource

	// Permission defines the access level of resource. It is based on both the
	// action and the subject who wants to perform the action.
	Permission(s Subject, action ...string) AccessLevel
}

// Subject instances are entities that want to perform action on Resource.
type Subject interface {
	// Relation returns the privilege value of the current Subject over passed
	// Subject.
	Relation(ctx any, s Subject) Relation
}

// Delegatee instances helps to limit the privileges of a Subject.
type Delegatee interface {
	// Delegate returns true if the condition is allowed to perform, and vice
	// versa.
	Delegate(relation Relation, resource Resource, action ...string) bool
}

// Checker supports check if a subject can perform action on resource or not.
type Checker struct {
	subject   Subject
	delegatee Delegatee
	action    []string
}

// Check returns a Checker with the subject.
func Check(s Subject) *Checker {
	return &Checker{subject: s}
}

// Delegate delegates the privileges to a delegatee.
func (c *Checker) Delegate(d Delegatee) *Checker {
	c.delegatee = d
	return c
}

// Perform adds action to Checker and returns itself.
func (c *Checker) Perform(action ...string) *Checker {
	c.action = action
	return c
}

// On checks if a subject can perform action on resource or not.
func (c *Checker) On(resource Resource) error {
	var accessLevel AccessLevel
	switch t := resource.(type) {
	case StaticResource:
		accessLevel = t.Permission(c.action...)
	case DynamicResource:
		accessLevel = t.Permission(c.subject, c.action...)
	default:
		panic(NotImplementedError.New(
			"expected an object implementing Resource or EnhancedResource"))
	}

	var ctx = resource.Context()
	var owner = resource.Owner()

	if ctx == owner && owner != nil {
		panic(XyprivError.New("do not use the owner as the context, you " +
			"should set the context as nil in this case"))
	}

	var privilege = Anyone
	if c.subject != nil {
		var relation = c.subject.Relation(ctx, owner)
		privilege = getPrivilege(ctx, relation)
		if c.delegatee != nil && !c.delegatee.Delegate(relation, resource, c.action...) {
			return PermissionError.Newf(
				"%s do not have the permission to %s %s",
				getName(c.subject), strings.Join(c.action, "_"), getName(resource))
		}
	}

	if int(privilege) < int(accessLevel) {
		return PermissionError.Newf(
			"%s do not have the permission to %s %s",
			getName(c.subject), strings.Join(c.action, "_"), getName(resource))
	}

	return nil
}

// getName returns the name of object.
func getName(a any) string {
	if a == nil {
		return ""
	}

	if s, ok := a.(fmt.Stringer); ok {
		return s.String()
	}

	var atype = reflect.TypeOf(a)
	switch atype.Kind() {
	case reflect.String:
		return a.(string)
	case reflect.Struct:
		return atype.Name()
	case reflect.Pointer:
		return atype.Elem().Name()
	default:
		panic(XyprivError.New("expected a string, struct, or pointer of struct to get its name"))
	}
}
