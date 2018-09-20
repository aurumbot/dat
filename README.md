# Bot Backend Tools

lib contains 2 modules, foundation and dat, that provides an array of
useful functions and standards for aurum.

## Foundation

Package foundation defines core bot information:
- The `Command` struct, its fields, and the `Action()` function which
  the bot uses to process commands.
- The `Bot` type, which provides fundamental configuration info such
  as authentication info and channels the bot should not respond on.
- A couple functions to dry up code
  - `GetGuild()` to get the server a message was sent in (as a 
    `discordgo.Guild` object)
  - `HasPermissions()` to check user's permissions against 
    [discord's built-in system](https://discordapp.com/developers/docs/topics/permissions)
  - `RoleFromID()` to get a `discordgo.Role` given its id as a string
- An exported version of `Bot` for plugins to call on
- An exported version of the `discordgo.Session` for plugins to call
  on

## dat

Package dat standardizes file I/O
- a `Log` of `log.Logger` which handles... logging errors.
- a `Save()` to write an `interface{}` to a file
- a `Load()` to retrieve an `interface{}` from a file.
- an `AlertDiscord()` which reads out an error to a given discord
  channel.

# EOF
