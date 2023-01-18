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

// LeastPrivilegeToken is a Token implements Delegatee. It uses the principle of
// least privilege.
type LeastPrivilegeToken struct {
	rules map[string]bool
}

// NewToken creates a LeastPrivilegeToken that implements Delegatee. It uses the
// principle of least privilege. By default, all privileges is rejected.
func NewToken() *LeastPrivilegeToken {
	return &LeastPrivilegeToken{
		rules: make(map[string]bool),
	}
}

// AllowAction allows all privileges on action.
func (t *LeastPrivilegeToken) AllowAction(action ...string) {
	t.setRule("", "", action, true)
}

// AllowScope allows all privileges in scope.
func (t *LeastPrivilegeToken) AllowScope(scope any) {
	t.setRule("", scope, nil, true)
}

// AllowRelation allows all privileges of relation in scope.
func (t *LeastPrivilegeToken) AllowRelation(relation Relation, scope any) {
	t.setRule(relation, scope, nil, true)
}

// AllowActionInScope allows all privileges on action in scope.
func (t *LeastPrivilegeToken) AllowActionInScope(scope any, action ...string) {
	t.setRule("", scope, action, true)
}

// Allow allows all privileges of relation on action in scope.
func (t *LeastPrivilegeToken) Allow(relation Relation, scope any, action ...string) {
	t.setRule(relation, scope, action, true)
}

// BanAction bans all privileges on action.
func (t *LeastPrivilegeToken) BanAction(action ...string) {
	t.setRule("", "", action, false)
}

// BanScope bans all privileges in scope.
func (t *LeastPrivilegeToken) BanScope(scope any) {
	t.setRule("", scope, nil, false)
}

// BanRelation bans all privileges of relation in scope.
func (t *LeastPrivilegeToken) BanRelation(relation Relation, scope any) {
	t.setRule(relation, scope, nil, false)
}

// BanActionInScope bans all privileges on action in scope.
func (t *LeastPrivilegeToken) BanActionInScope(scope any, action ...string) {
	t.setRule("", scope, action, false)
}

// Ban bans all privileges of relation on action in scope.
func (t *LeastPrivilegeToken) Ban(relation Relation, scope any, action ...string) {
	t.setRule(relation, scope, action, false)
}

// Delegate checks if condition tuple is allowed or banned.
func (t LeastPrivilegeToken) Delegate(relation Relation, resource Resource, action ...string) bool {
	var relName = string(relation)
	var rsrName = getName(resource)
	var ctxName = getName(resource.Context())
	var actName = strings.Join(action, "_")

	var keys = []string{
		// Full-parameter keys.
		strings.Join([]string{actName, relName, rsrName}, "."),
		strings.Join([]string{actName, relName, ctxName}, "."),

		// Partial keys.
		strings.Join([]string{"", relName, rsrName}, "."),
		strings.Join([]string{"", relName, ctxName}, "."),
		strings.Join([]string{actName, "", rsrName}, "."),
		strings.Join([]string{actName, "", ctxName}, "."),

		// One-parameter keys.
		strings.Join([]string{"", "", rsrName}, "."),
		strings.Join([]string{"", "", ctxName}, "."),
		strings.Join([]string{actName, "", ""}, "."),
	}

	var isAllow = false
	for _, k := range keys {
		if val, ok := t.rules[k]; ok {
			if !val {
				return false
			}
			isAllow = true
		}
	}

	return isAllow
}

// setRule adds the condition tuple of relation, scope, and action into token
// rules. The scope could be context or resource. Use the empty relation to
// apply all relations in the condition. Use the empty string as scope to apply
// all scopes in the condition. Use no action to apply all actions in the
// condition.
func (t *LeastPrivilegeToken) setRule(relation Relation, scope any, action []string, result bool) {
	var relName = string(relation)
	var scopeName = getName(scope)
	var actName = strings.Join(action, "_")

	t.rules[strings.Join([]string{actName, relName, scopeName}, ".")] = result
}
