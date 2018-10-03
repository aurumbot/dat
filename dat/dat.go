package dat

import (
	"bytes"
	"encoding/json"
	dsg "github.com/bwmarrin/discordgo"
	"io"
	"log"
	"os"
	"strings"
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

// Sets the absolute path
func SetPath(p string) {
	path = p + "/dat/"
}

// Simple check to make sure a... something... exists.
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

var lock sync.Mutex

/* Saves a file to json
* This is a simple helper function to manage data persistance and keep it in
* one place.
*
* Parameters
* - fileName (string) : where the data is to be stored. Note that this should
*			be in the form of $MODULENAME/*.json so modules that
*			use the same name config.json don't interfere.
* - v (interface{})   : the thing you're saving.
*
* Returns:
* - error : error has already been logged, but useful for AlertDiscord()
 */
func Save(fileName string, item interface{}) error {
	lock.Lock()
	defer lock.Unlock()
	// Checks for directories in fileName, creates them if they aren't real.
	dir := ""
	split := strings.Split(fileName, "/")
	for i, v := range split {
		dir += v
		if len(split)-1 != i {
			dir += "/"
		} else {
			break
		}
		// makes sure its real
		validPath, err := exists(dir)
		if err != nil {
			Log.Println(err)
			return err
		}
		if validPath {
			continue
		}
		// creates
		if err := os.Mkdir(dir, 0600); err != nil {
			Log.Println(err)
			return err
		}
	}
	// Opens file, creates if not real
	file, err := os.OpenFile(path+"cfg/"+fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		Log.Println(err)
		return err
	}
	defer file.Close()

	b, err := json.MarshalIndent(item, "", "\t")
	if err != nil {
		Log.Println(err)
		return err
	}
	reader := bytes.NewReader(b)
	_, err = io.Copy(file, reader)
	return err
}

/* Loads a file from json
* This is a simple helper function to manage data persistance and keep it in
* one place.
*
* Parameters
* - fileName (string) : where the data is stored. Note that this should
*			be in the form of $MODULENAME/*.json so modules that
*			use the same name config.json don't interfere.
* - v (interface{})   : the thing you're loading. as a pointer
*
* Returns:
* - error : error has already been logged, but useful for AlertDiscord()
 */
func Load(fileName string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()
	file, err := os.OpenFile(path+"cfg/"+fileName, os.O_RDONLY, 0600)
	if err != nil {
		Log.Println(err)
		return err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(v)
	return err
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
