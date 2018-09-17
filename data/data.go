package dat

import (
	"bytes"
	"encoding/json"
	f "github.com/aurumbot/dat/foundation"
	dsg "github.com/bwmarrin/discordgo"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

var (
	currentTime string
	Log         *log.Logger
	path        = "./dat/"
)

func init() {
	currentTime = time.Now().Format("2006-01-02@15h04m")

	file, err := os.Create(path + "logs/botlogs@" + currentTime + ".log")
	if err != nil {
		panic(err)
	}
	Log = log.New(file, "", log.Ldate|log.Ltime|log.Llongfile|log.LUTC)
}

var lock sync.Mutex

func Save(fileName string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()

	file, err := os.Create(path + "cmds/" + fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	reader := bytes.NewReader(b)

	_, err = io.Copy(file, reader)

	return err
}

func Load(fileName string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()
	file, err := os.Open(path + "cmds/" + fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(v)
	return err
}

// TODO: delete, rewrite main to pull from a Load() call.
func GetBotInfo() (f.BotType, error) {
	raw, err := ioutil.ReadFile("./dat/preferences.json")
	var b f.BotType

	if err != nil {
		return b, err
	}

	err = json.Unmarshal(raw, &b)

	if err != nil {
		return b, err
	}

	return b, nil
}

/* # Alerts discord of errors.
* AlertDiscord is a function that... well alerts discord if there's a problem.
* Useful for things like if your command fails and you have to return, the user
* isn't kept in limbo waiting for something to happen. However this is not a
* substitute for posting an error in the log and should be done *along with*
* dat.Log.New(), this just helps prevent the users moaning about "broken bot"
* and actually proves it to them.
*
* Parameters:
* - s (type *discordgo.Session) : Needed for posting a message
* - m (type *discordgo.Message) : Needed for posting a message. Pings .Author.
* - err (type error) : The error being reported
 */
func AlertDiscord(s *dsg.Session, m *dsg.Message, err error) {
	str := `<@` + m.Author.ID + `> | Error encountered, details as follows:
	` + "\n```" + err.Error() + "```\n" + `
You are being pinged because your message was the message that triggered the 
above error. Please inform the person running this bot or a sever admin.`
	s.ChannelMessageSend(m.ChannelID, str)
}
