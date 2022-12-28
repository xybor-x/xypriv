# Subjects

-   Group

-   User

# Context and relations

## Non-context

-   User:
    -   self
    -   admin
    -   moderator
    -   anyone

-   Group:
    -   anyone

## Normal context

-   Use the recommended privileges and access levels.

-   User:
    -   friend (low familiar): with another user.
    -   banned (bad relation): with another user.
    -   group admin (local admin): with group.
    -   group member (low familiar): with group.
    -   group banned (bad relation): with group.

## Group context

-   Use the recommended privileges and access levels.

-   User:
    -   anyone: with user (if one of users doesn't join in group).
    -   banned (bad relation): with user (if one of users is banned by group, or users banned each other).
    -   group admin (local admin): with user.
    -   same group (low familiar): with user.

# Resources

## Avatar

-   Owner: User
-   Context: The normal context.
-   Actions:
    -   update: admin.
    -   comment: friend.
    -   read: anyone.

## GroupPost

-   Owner: User
-   Context: The group context.
-   Actions:
    -   update: self.
    -   delete: group admin.
    -   read: same group.

## GroupAvatar

-   Owner: Group
-   Context: The normal context.
-   Actions:
    -   update: group admin.
    -   delete: group admin.
    -   read: anyone.
