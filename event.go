package main

import (
	"fmt"
	"regexp"
	"io/ioutil"

	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

// You more than likely want your "Bot User OAuth Access Token"
func getOauthToken() string {
	dat, err := ioutil.ReadFile("./secrets/oauth")
	if (err != nil) { panic(err) }
	var ret = string(dat)
	return ret[:len(ret)-1]
}

func init() {
	var oauthToken = getOauthToken();
	var api = slack.New(oauthToken);
	api.PostMessage("CCUDKKFT5", slack.MsgOptionText("_*info* - I have restarted_", false))
}

func checkWeightMeasurement(message string) {
	var oauthToken = getOauthToken();
	var api = slack.New(oauthToken);
	var re, err = regexp.MatchString("^\\d+\\.\\d+ this morning$", message)
	if (err != nil) {
		api.PostMessage("CCUDKKFT5", slack.MsgOptionText("This is too much to handle. I quit.", false))
		panic(err)
	}
	if (re) {
		api.PostMessage("CCUDKKFT5", slack.MsgOptionText("I received a weight measurement", false))
	} else {
		api.PostMessage("CCUDKKFT5", slack.MsgOptionText("I did not receive a weight measurement", false))
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
		if (ev.Channel == "CCUDKKFT5" && ev.BotID == "") {
			checkWeightMeasurement(ev.Text)
			// api.PostMessage(ev.Channel, slack.MsgOptionText(ev.Text, false))
		}
	}
}

