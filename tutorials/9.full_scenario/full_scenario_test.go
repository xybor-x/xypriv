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
	"fmt"
	"reflect"

	"github.com/xybor-x/xypriv"
)

func init() {
	xypriv.AddRelation(nil, "banned", xypriv.BadRelation)
	xypriv.AddRelation(nil, "friend", xypriv.LowFamiliar)
	xypriv.AddRelation(nil, "groupMember", xypriv.LowFamiliar)
	xypriv.AddRelation(nil, "groupAdmin", xypriv.LocalAdmin)
	xypriv.AddRelation(nil, "groupBanned", xypriv.BadRelation)
	xypriv.AddRelation(Group{}, "banned", xypriv.BadRelation)
	xypriv.AddRelation(Group{}, "anyone", xypriv.Anyone)
	xypriv.AddRelation(Group{}, "sameGroup", xypriv.LowFamiliar)
	xypriv.AddRelation(Group{}, "moderator", xypriv.Moderator)
	xypriv.AddRelation(Group{}, "groupAdmin", xypriv.LocalAdmin)
	xypriv.AddRelation(Group{}, "admin", xypriv.Admin)
	xypriv.AddRelation(Group{}, "self", xypriv.Self)
}

func Example() {
	var self = User{id: "self"}
	var admin = User{id: "admin", role: "admin"}
	var gAdmin = User{id: "gAdmin"}
	var sameGroup = User{id: "sameGroup"}
	var bannedButSameGroup = User{id: "bannedButSameGroup"}
	var friendButNotSameGroup = User{id: "friendButNotSameGroup"}
	var userBannedByGroup = User{id: "userBannedByGroup"}

	self.bannedUser = append(self.bannedUser, bannedButSameGroup)
	self.friends = append(self.friends, friendButNotSameGroup)

	var groupX = Group{
		members:    []User{self, gAdmin, sameGroup, bannedButSameGroup},
		admin:      gAdmin,
		bannedUser: []User{userBannedByGroup},
	}

	var avt = Avatar{user: self}
	var groupAvt = GroupAvatar{group: groupX}
	var post = GroupPost{group: groupX, user: self}

	var users = []User{
		self,
		admin,
		gAdmin,
		sameGroup,
		bannedButSameGroup,
		friendButNotSameGroup,
		userBannedByGroup,
	}

	var resources = []struct {
		r       xypriv.Resource
		actions []string
	}{
		{avt, []string{"update", "read", "comment"}},
		{groupAvt, []string{"update", "delete", "read"}},
		{post, []string{"update", "delete", "read"}},
	}

	for _, resource := range resources {
		var rname = reflect.TypeOf(resource.r).Name()
		fmt.Printf("===================== %s ======================\n", rname)
		for _, a := range resource.actions {
			for _, u := range users {
				var uname = u.id
				if xypriv.Check(u).Perform(a).On(resource.r) == nil {
					fmt.Printf("%s CAN %s %s\n", uname, a, rname)
				} else {
					fmt.Printf("%s CAN NOT %s %s\n", uname, a, rname)
				}
			}
			fmt.Println("-----------------------------------------------")
		}
	}

	// Output:
	// ===================== Avatar ======================
	// self CAN update Avatar
	// admin CAN update Avatar
	// gAdmin CAN NOT update Avatar
	// sameGroup CAN NOT update Avatar
	// bannedButSameGroup CAN NOT update Avatar
	// friendButNotSameGroup CAN NOT update Avatar
	// userBannedByGroup CAN NOT update Avatar
	// -----------------------------------------------
	// self CAN read Avatar
	// admin CAN read Avatar
	// gAdmin CAN read Avatar
	// sameGroup CAN read Avatar
	// bannedButSameGroup CAN NOT read Avatar
	// friendButNotSameGroup CAN read Avatar
	// userBannedByGroup CAN read Avatar
	// -----------------------------------------------
	// self CAN comment Avatar
	// admin CAN comment Avatar
	// gAdmin CAN NOT comment Avatar
	// sameGroup CAN NOT comment Avatar
	// bannedButSameGroup CAN NOT comment Avatar
	// friendButNotSameGroup CAN comment Avatar
	// userBannedByGroup CAN NOT comment Avatar
	// -----------------------------------------------
	// ===================== GroupAvatar ======================
	// self CAN NOT update GroupAvatar
	// admin CAN update GroupAvatar
	// gAdmin CAN update GroupAvatar
	// sameGroup CAN NOT update GroupAvatar
	// bannedButSameGroup CAN NOT update GroupAvatar
	// friendButNotSameGroup CAN NOT update GroupAvatar
	// userBannedByGroup CAN NOT update GroupAvatar
	// -----------------------------------------------
	// self CAN NOT delete GroupAvatar
	// admin CAN delete GroupAvatar
	// gAdmin CAN delete GroupAvatar
	// sameGroup CAN NOT delete GroupAvatar
	// bannedButSameGroup CAN NOT delete GroupAvatar
	// friendButNotSameGroup CAN NOT delete GroupAvatar
	// userBannedByGroup CAN NOT delete GroupAvatar
	// -----------------------------------------------
	// self CAN read GroupAvatar
	// admin CAN read GroupAvatar
	// gAdmin CAN read GroupAvatar
	// sameGroup CAN read GroupAvatar
	// bannedButSameGroup CAN read GroupAvatar
	// friendButNotSameGroup CAN read GroupAvatar
	// userBannedByGroup CAN NOT read GroupAvatar
	// -----------------------------------------------
	// ===================== GroupPost ======================
	// self CAN update GroupPost
	// admin CAN NOT update GroupPost
	// gAdmin CAN NOT update GroupPost
	// sameGroup CAN NOT update GroupPost
	// bannedButSameGroup CAN NOT update GroupPost
	// friendButNotSameGroup CAN NOT update GroupPost
	// userBannedByGroup CAN NOT update GroupPost
	// -----------------------------------------------
	// self CAN delete GroupPost
	// admin CAN delete GroupPost
	// gAdmin CAN delete GroupPost
	// sameGroup CAN NOT delete GroupPost
	// bannedButSameGroup CAN NOT delete GroupPost
	// friendButNotSameGroup CAN NOT delete GroupPost
	// userBannedByGroup CAN NOT delete GroupPost
	// -----------------------------------------------
	// self CAN read GroupPost
	// admin CAN read GroupPost
	// gAdmin CAN read GroupPost
	// sameGroup CAN read GroupPost
	// bannedButSameGroup CAN NOT read GroupPost
	// friendButNotSameGroup CAN NOT read GroupPost
	// userBannedByGroup CAN NOT read GroupPost
	// -----------------------------------------------
}
