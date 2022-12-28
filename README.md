[![xybor founder](https://img.shields.io/badge/xybor-huykingsofm-red)](https://github.com/huykingsofm)
[![Go Reference](https://pkg.go.dev/badge/github.com/xybor-x/xypriv.svg)](https://pkg.go.dev/github.com/xybor-x/xypriv)
[![GitHub Repo stars](https://img.shields.io/github/stars/xybor-x/xypriv?color=yellow)](https://github.com/xybor-x/xypriv)
[![GitHub top language](https://img.shields.io/github/languages/top/xybor-x/xypriv?color=lightblue)](https://go.dev/)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/xybor-x/xypriv)](https://go.dev/blog/go1.18)
[![GitHub release (release name instead of tag name)](https://img.shields.io/github/v/release/xybor-x/xypriv?include_prereleases)](https://github.com/xybor-x/xypriv/releases/latest)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/8060ee5daef141c6811223384c8da55f)](https://www.codacy.com/gh/xybor-x/xypriv/dashboard?utm_source=github.com&utm_medium=referral&utm_content=xybor-x/xypriv&utm_campaign=Badge_Grade)
[![Go Report](https://goreportcard.com/badge/github.com/xybor-x/xypriv)](https://goreportcard.com/report/github.com/xybor-x/xypriv)

# Introduction

Package xypriv supports to central manage and control privileges that protect
resources from illegal access.

# Definitions

**Resources** are objects that need to be protected from illegal _actions_.

**Actions** are activities which are performed on _resources_.

**Subjects** are entities who want to perform _actions_ on _resources_.

There are some **Subjects** are considered as the owner of **Resources**.

# How it works

The accessible ability of a subject on resource is based on:

1.  The relation between subject and owner. It defines a privilege number.

    _The higher the **privilege**, the **more easily** the subject can access on the owner's resources_.

2.  The permission set of resource, it contains many elements, each element is
    a tuple of an owner and an access level.

    _The higher the **access level**, the **more difficult** a subject can access on the owner's resources_.

Assuming that:

1.  The resource has the owner, called `O`.
2.  The resource defines the access level to read is `AC`.
3.  The subject has the privilege over `O` is `P`.

The subject can read the resource if `P >= AC`.

# Get started

_After you read the following guide, please visit [tutorials](./tutorials/)._

## Define entities

First, you need to define all subjects, resources, contexts, and relations in
each context. For example, see [tutorials](./tutorials/).

Then, give prefered names to the privileges and access levels in each context. Below items are recommended privileges and access levels for the normal context.

| Privilege |       Relation       |
| :-------- | :------------------: |
| 10        |         Self         |
| 8-9       |         Admin        |
| 6-7       | Moderator, Supporter |
| 2-5       |       Familiars      |
| 1         |        Anyone        |
| 0         |     Bad relation     |

| Permission |      Level     |
| :--------- | :------------: |
| 10         |   Top secret   |
| 8-9        |     Secret     |
| 6-7        |  Confidential  |
| 2-5        |     Private    |
| 1          |     Public     |
| ~~0~~      | ~~Do not use~~ |

## Implement Subject interface

_After you read the all sections, please visit [tutorials](./tutorials/)._

`Subject` interface requires only one method called `Relation`.

It receives another `Subject` as parameter and the context. It expects to return
the _privilege_ number of the current `Subject` over the parameter `Subject`.

The latter parameter helps to find relations depending on the context. The
context can be seen as the namespace of relations.

For example, assuming that user A joined in group G, user B is the admin of
group G, user C is a friend of user A:

1.  In the normal context, user B has no relation with user A, user C is the
    friend of user A.
2.  In the context of group G, user B is the group admin of user A, user C has
    no relation with user A.

The context will be determined by the Resource. The normal context should be
represented by `nil`.

```golang
type Subject interface {
    Relation(s Subject, ctx any) Privilege
}
```

For example:

```golang
// User implements Subject interface.
type User struct {
    id string
}

func (u User) Relation(s Subject, ctx any) xypriv.Privilege {
    switch t := s.(type) {
    case User:
        if u.id == t.id {
            return xypriv.Self
        }
    }
    switch c := ctx.(type) {
    case nil:
        // handle context-dependent privileges.
    }
    return xypriv.Anyone
}
```

## Implement Resource interfaces

_After you read the all sections, please visit [tutorials](./tutorials/)._

For a resource, you need to define `Context` and `Owner` methods which returns
the context and owner of resource, respectively.

You also need to define `Permission` method, it returns the `AccessLevel` to
perform the action on the resource. The parameter of this method depends on
which problem you are solving.

```golang
type Resource interface {
    Context() any
    Owner() Subject
}

type SimpleResource interface {
    Resource
    Permission(action ...string) AccessLevel
}

type EnhancedResource interface {
    Resource
    Permission(s Subject, action ...string) AccessLevel
}
```

For example:

```golang
// Avatar implements SimpleResource.
type Avatar struct {
    user User
}

func (a Avatar) Context() any {
    return nil
}

fuunc (a Avatar) Owner() xypriv.Subject {
    return a.user
}

func (a Avatar) Permission(action ...string) xypriv.AccessLevel {
	if len(action) != 1 {
		return xypriv.NotSupport
	}

	switch action[0] {
	case "update":
		return xypriv.LowPrivate
	}
	return xypriv.NotSupport
}
```

## Use Check-Perform-On operator

_After you read the all sections, please visit [tutorials](./tutorials/)._

```golang
var userA = User{name: "A"}
var userB = User{name: "B"}

var avtA = Avatar{user: userA}

xypriv.Check(userA).Perform("update").On(avtA) // nil
xypriv.Check(userB).Perform("update").On(avtA) // XyprivError
```
