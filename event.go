package main

import (
	"fmt"
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
			api.PostMessage(ev.Channel, slack.MsgOptionText(ev.Text, false))
		}
	}
}

