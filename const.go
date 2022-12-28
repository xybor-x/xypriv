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
	"strings"
)

// Recommended privileges.
const (
	BadRelation Privilege = iota
	Anyone
	LowFamiliar
	MediumFamiliar
	HighFamiliar
	TopFamiliar
	LocalModerator
	Moderator
	LocalAdmin
	Admin
	Self
)

// Recommended access levels.
const (
	Public AccessLevel = iota + 1
	LowPrivate
	MediumPrivate
	HighPrivate
	TopPrivate
	LowConfidential
	HighConfidential
	LowSecret
	HighSecret
	TopSecret
	NotSupport
)

// defaultRelation stored all relations for all context.
var defaultRelation = map[Relation]Privilege{
	"badrelation":    BadRelation,
	"anyone":         Anyone,
	"lowfamiliar":    LowFamiliar,
	"mediumfamiliar": MediumFamiliar,
	"highfamiliar":   HighFamiliar,
	"topfamiliar":    TopFamiliar,
	"localmoderator": LocalModerator,
	"moderator":      Moderator,
	"localadmin":     LocalAdmin,
	"admin":          Admin,
	"self":           Self,
}

// relationMap stores all contexts and their relations.
var relationMap = map[string]map[Relation]Privilege{}

// AddRelation adds a relation of context to privilege manager. The context
// should be a string, struct, or pointer of struct.
func AddRelation(context any, relation Relation, privilege Privilege) {
	var cname = getName(context)

	if _, ok := relationMap[cname]; !ok {
		relationMap[cname] = make(map[Relation]Privilege)
	}

	relation = Relation(strings.ToLower(string(relation)))
	relationMap[cname][relation] = privilege
}

// getPrivilege returns the privilege corresponding to context and relation.
func getPrivilege(context any, relation Relation) Privilege {
	var cname = getName(context)

	if cmap, ok := relationMap[cname]; ok || cname == "" {
		relation = Relation(strings.ToLower(string(relation)))
		if priv, ok := cmap[relation]; ok {
			return priv
		}

		if priv, ok := defaultRelation[relation]; ok {
			return priv
		}

		panic(XyprivError.Newf("unknown relation %s in context %s", relation, cname))
	}
	panic(XyprivError.Newf("unknown context %s", cname))
}
