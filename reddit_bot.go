// main file
package main

import (
	"os"

	"github.com/mattermost/mattermost-server/model"
)

type user struct {
	Mail     string
	Name     string
	Password string
	First    string
	Last     string
	Team     string
}

type config struct {
	User     user
	channels map[string]*model.Channel
}

var client *model.Client4
var botUser *model.User
var currentConfig *config

func main() {
	client = model.NewAPIv4Client("http://localhost:8065")

	CheckServerStatus()
	LoadConfig()
	SendMessageToChannel("debugging-for-sample-bot", "_Starting bot..._", "")
}

// this is supposed to load config from a file or database idk.
// fills mock data for testing atm
func LoadConfig() {
	cconfig := config{}
	currentConfig = &cconfig

	cuser := user{}
	cuser.Mail = "bot@example.com"
	cuser.Password = "password1"
	cuser.Name = "Reddit Bot"
	cuser.First = "Reddit"
	cuser.Last = "Bot"
	cuser.Team = "test"

	cconfig.User = cuser

	Login()

	var team *model.Team
	var resp *model.Response
	if team, resp = client.GetTeamByName(cuser.Team, ""); resp.Error != nil {
		LogError(resp.Error)
		os.Exit(1)
	}

	CHANNEL_LOG_NAME := "debugging-for-sample-bot"

	currentConfig.channels = make(map[string]*model.Channel, 0)
	var rchannel *model.Channel
	if rchannel, resp = client.GetChannelByName(CHANNEL_LOG_NAME, team.Id, ""); resp.Error != nil {
		LogError(resp.Error)
	}

	currentConfig.channels[CHANNEL_LOG_NAME] = rchannel
}

func CheckServerStatus() {
	if _, resp := client.GetOldClientConfig(""); resp.Error != nil {
		LogError(resp.Error)
		os.Exit(1)
	}
}

func Login() {
	tempUsr := (*currentConfig).User

	if user, resp := client.Login(tempUsr.Mail, tempUsr.Password); resp.Error != nil {
		LogError(resp.Error)
		os.Exit(1)
	} else {
		botUser = user
	}
}

func SendMessageToChannel(channel string, msg string, replyToId string) {
	post := &model.Post{}

	chann := currentConfig.GetChannel(channel)
	post.ChannelId = chann.Id
	post.Message = msg

	post.RootId = replyToId

	if _, resp := client.CreatePost(post); resp.Error != nil {
		LogError(resp.Error)
	}
}

func (c *config) GetChannel(id string) *model.Channel {
	return (*c).channels[id]
}

func LogError(err *model.AppError) {
	println("Error: ", err.Message)
	println(err.DetailedError)
}
