package main

import (
	"os"
	"fmt"
	"time"
	"regexp"
	"strconv"
	"io/ioutil"

	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

// You more than likely want your "Bot User OAuth Access Token"
func getOauthToken() string {
	dat, err := ioutil.ReadFile("./secrets/oauth")
	if (err != nil) {
		panic(err)
	}
	var ret = string(dat)
	return ret[:len(ret)-1]
}

func init() {
	var oauthToken = getOauthToken();
	var api = slack.New(oauthToken);
	api.PostMessage("CCUDKKFT5", slack.MsgOptionText("_*info* - I have restarted_", false))
}

func checkWeightMeasurement(message string, user string, channel string) {
	var oauthToken = getOauthToken();
	var api = slack.New(oauthToken);
	var re = regexp.MustCompile("^(\\d+\\.\\d+) this morning$")
	var match = re.FindStringSubmatch(message)
	if (len(match) >= 2) {
		// TODO: store results in a google spreadsheet
		var dir = "./data/" + user
		os.Mkdir(dir, 0640)
		var weightFile = dir + "/weight"
		f, err := os.OpenFile(weightFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
		if (err != nil) {
			panic(err)
		}
		if _, err := f.Write([]byte(strconv.FormatInt(time.Now().Unix(), 10) + "," + match[1] + "\n")); err != nil {
			panic(err)
		}
		if err := f.Close(); err != nil {
			panic(err)
		}
		data, err := ioutil.ReadFile(weightFile)
		if (err != nil) {
			panic(err)
		}
		// TODO: parse, and form a more useful response:
		api.PostMessage(channel, slack.MsgOptionText(string(data), false))
	}
}

func handle(event slackevents.EventsAPIInnerEvent) {
	var oauthToken = getOauthToken();
	var api = slack.New(oauthToken);
	fmt.Printf("event: %v", event)
	switch ev := event.Data.(type) {
	case *slackevents.AppMentionEvent:
		api.PostMessage(ev.Channel, slack.MsgOptionText("At your service.", false))
	case *slackevents.MessageEvent:
		fmt.Printf("ev: %v", ev)
		if (ev.Channel == "CATMEMUFK") {
			// #training
			checkWeightMeasurement(ev.Text, ev.User, ev.Channel)
		}
		if (ev.Channel == "CCUDKKFT5" && ev.BotID == "") {
			if (ev.Text == "good bot") {
				api.PostMessage(ev.Channel, slack.MsgOptionText("thank you, human", false))
			}
			// #bot_school
			// api.PostMessage(ev.Channel, slack.MsgOptionText(ev.Text, false))
		}
	}
}

