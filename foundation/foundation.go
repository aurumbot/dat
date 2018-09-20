package foundation

import (
	dsg "github.com/bwmarrin/discordgo"
)

/* # Big Bot Boye
* This struct defines the authentication info and config details for the bot
*
* The fields are self explanitory so I'm not going to detail them.
*
* TODO: I think the auth section needs to be either a different json file or
*       passed in as parameters because the rest isn't nearly as important.
*
 */
type Bot struct {
	ClientID            string   `json:"clientID"`
	Secret              string   `json:"secret"`
	Token               string   `json:"token"`
	Prefix              string   `json:"prefix"`
	BlacklistedChannels []string `json:"channels"`
	BlacklistedRoles    []string `json:"roles"`
	Admins              []string `json:"admins"`
}

/* Defines the actual action the bot takes
* This is a key component of the Command struct, it is the actual *thing* that
* the command does when it is run. This can be as simple as printing a string
* or embed to discord or even be a wrapper for a massive series of functions
* for your discord-based RPG.
*
* Parameters:
* - session (*discordgo.Session) The bot session, in case you need to pull data
*	from about discord itself to complete your task
* - session (*discordgo.Message) The entire message that triggered the command
*	this includes the prefix and command itself. You will have to parse out
*	flags and clean the input.
*
* NOTE: THIS DOES NOT RETURN ERRORS. YOU MUST HANDLE ERRORS.
 */
type Action func(session *dsg.Session, message *dsg.Message)

/* Defines static data about commands the bot runs.
* This is a very large structure that defines all the needed bits for a bot
* command. All bot modules MUST have one of these along with a few other key
* components so that the bot works.
*
* Parameters:
* - Name (string)   | The name of the command, used in the help page
* - Help (string)   | A description of the command and how it is used.
* - Action (Action) | The function the command runs
* - Perms (Int)     | Oh boy this is a doozie:
* This is the integer value of the permissions needed to run the command. If
* you don't know what value to use, use a calculator tool such as
* https://discordapi.com/permissions.html or discordgo's predefined constants
* available at https://godoc.org/github.com/bwmarrin/discordgo#pkg-constants
* If you want your command available to all users who can trigger the bot, set
* this value to -1.
 */
type Command struct {
	Name    string `json:"name"`
	Help    string `json:"help"`
	Perms   int    `json:"perms"`
	Version string `json:"version"`
	Action  Action `json:"-"`
}

/* # Get the guild a message was sent in.
* What a pain in the arse.
*
* Parameters:
* - s (type *discordgo.Session) | The current running discord session,
*     (discordgo needs that always apparently)
* - message (type *discordgo.Message) | the author's id is extracted from this.
*
* Returns:
* - st (type *discordgo.Guild) | The guild the message was found in
* - err (type error)           | If an error was encountered during the process
*	This error is an SEP (someone else's problem).
 */
func GetGuild(s *dsg.Session, m *dsg.Message) (st *dsg.Guild, err error) {
	chn, err := s.Channel(m.ChannelID)
	if err != nil {
		return &dsg.Guild{}, err
	}

	gid := chn.GuildID

	return s.Guild(gid)
}

/* Checks if user has permission to run a command
* This function is a wrapper to check if a user has the permission needed to
* run a given command. This checks for both specific permissions the user has
* in the server (see below) and for "bot staff" roles defined in the config.
* Permissions are integer constants defined by discordgo:
* https://godoc.org/github.com/bwmarrin/discordgo#pkg-constants
* Note that the check is non-hierarchichal.
*
* Parameters:
* - s (*discordgo.Session)
* - m (*discordgo.Message)
* - userID (string) : The user to be checked, leave blank for the message
*			author.
* - perm (int) : see above paragraph
*
* Returns:
* - bool  : if the perm is met
* - error : if an error was encountered, must be logged.
 */
func HasPermissions(s *dsg.Session, m *dsg.Message, userID string, perm int) (bool, error) {
	guild, err := GetGuild(s, m)
	if err != nil {
		return false, err
	}
	var member *dsg.Member
	if userID == "" {
		member, err = s.GuildMember(guild.ID, m.Author.ID)
	} else {
		member, err = s.GuildMember(guild.ID, userID)
	}
	if err != nil {
		return false, err
	}
	for _, b := range Config.Admins {
		if Contains(member.Roles, b) {
			return true, nil
		}
	}
	for _, b := range member.Roles {
		role, err := RoleFromID(s, m, b)
		if err != nil {
			return false, err
		}

		if role.Permissions&perm != 0 {
			return true, nil
		} else if role.Permissions&dsg.PermissionAdministrator != 0 {
			return true, nil
		}
	}
	return false, nil
}

// a stupid, inefficent function to get a role from its id
func RoleFromID(s *dsg.Session, m *dsg.Message, id string) (*dsg.Role, error) {
	guild, err := GetGuild(s, m)
	if err != nil {
		return &dsg.Role{}, err
	}
	roles, err := s.GuildRoles(guild.ID)
	if err != nil {
		return &dsg.Role{}, err
	}

	for _, role := range roles {
		if role.ID == id {
			return role, nil
		}
	}
	return &dsg.Role{}, nil
}

/* # Check if item is in array
* This function checks if a value is in a slice (string only)
*
* Parameters:
* - list ([]string) | the slice to be checking against
* - item (string)   | the item looked for in the slice
*
* Returns:
* - bool | If the item was found or not
*
* NOTE: If another Contains() funciton is needed for a different type, rename
* this function to ContainsSliceString() and the other function to
* ContainsSlice<T>() where <T> is the generic type.
 */
func Contains(list []string, item string) bool {
	for _, b := range list {
		if b == item {
			return true
		}
	}
	return false
}

// An initialized instance of the Bot for use everywhere in this project.
var Config Bot

// A discordgo session global variable as .Session is needed for a lot.
// This is written to in the main.go file.
var Session *dsg.Session
